package bilibili

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WBI (Web Browser Interface) 签名管理器
// 用于B站API请求的安全签名验证
type WBIManager struct {
	imgKey     string
	subKey     string
	mixinKey   string
	lastUpdate time.Time
	mutex      sync.RWMutex
}

// WBIKeys WBI密钥信息
type WBIKeys struct {
	ImgURL string `json:"img_url"`
	SubURL string `json:"sub_url"`
}

// WBI相关常量
const (
	WBIUpdateInterval = 2 * time.Hour // 密钥更新间隔
	WBIKeyLength      = 32            // 混合密钥长度
)

// WBI密钥映射表 (参考biliup实现)
var wbiKeyMap = [64]int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35,
	27, 43, 5, 49, 33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13,
	37, 48, 7, 16, 24, 55, 40, 61, 26, 17, 0, 1, 60, 51, 30, 4,
	22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11, 36, 20, 34, 44, 52,
}

// NewWBIManager 创建WBI管理器
func NewWBIManager() *WBIManager {
	return &WBIManager{}
}

// extractKey 从URL中提取key
func extractKey(urlStr string) string {
	if urlStr == "" {
		return ""
	}
	
	// 找到最后一个 '/' 和第一个 '.'
	lastSlash := strings.LastIndex(urlStr, "/")
	if lastSlash == -1 {
		return ""
	}
	
	firstDot := strings.Index(urlStr[lastSlash:], ".")
	if firstDot == -1 {
		return urlStr[lastSlash+1:]
	}
	
	return urlStr[lastSlash+1 : lastSlash+firstDot]
}

// UpdateKeys 更新WBI密钥
func (w *WBIManager) UpdateKeys(imgURL, subURL string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	imgKey := extractKey(imgURL)
	subKey := extractKey(subURL)
	
	if imgKey == "" || subKey == "" {
		return fmt.Errorf("failed to extract keys from URLs: img=%s, sub=%s", imgURL, subURL)
	}
	
	w.imgKey = imgKey
	w.subKey = subKey
	w.mixinKey = w.generateMixinKey(imgKey, subKey)
	w.lastUpdate = time.Now()
	
	return nil
}

// generateMixinKey 生成混合密钥
func (w *WBIManager) generateMixinKey(imgKey, subKey string) string {
	full := imgKey + subKey
	key := make([]byte, WBIKeyLength)
	
	for i := 0; i < WBIKeyLength; i++ {
		if wbiKeyMap[i] < len(full) {
			key[i] = full[wbiKeyMap[i]]
		}
	}
	
	return string(key)
}

// IsExpired 检查密钥是否过期
func (w *WBIManager) IsExpired() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	
	return time.Since(w.lastUpdate) >= WBIUpdateInterval
}

// Sign 为请求参数生成WBI签名
func (w *WBIManager) Sign(params map[string]string) map[string]string {
	w.mutex.RLock()
	mixinKey := w.mixinKey
	w.mutex.RUnlock()
	
	if mixinKey == "" {
		// 如果没有密钥，返回原参数（不签名）
		return params
	}
	
	// 复制参数避免修改原始数据
	signedParams := make(map[string]string)
	for k, v := range params {
		// 过滤特殊字符
		filtered := strings.Map(func(r rune) rune {
			switch r {
			case '!', '\'', '(', ')', '*':
				return -1 // 删除这些字符
			default:
				return r
			}
		}, v)
		signedParams[k] = filtered
	}
	
	// 添加时间戳
	timestamp := time.Now().Unix()
	signedParams["wts"] = strconv.FormatInt(timestamp, 10)
	
	// 对参数进行排序
	keys := make([]string, 0, len(signedParams))
	for k := range signedParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	// 构建查询字符串
	values := url.Values{}
	for _, k := range keys {
		values.Set(k, signedParams[k])
	}
	queryString := values.Encode()
	
	// 生成签名
	h := md5.New()
	h.Write([]byte(queryString + mixinKey))
	sign := fmt.Sprintf("%x", h.Sum(nil))
	
	// 添加签名到参数中
	signedParams["w_rid"] = sign
	
	return signedParams
}

// GetWBIKeys 获取WBI密钥信息
func (c *Client) GetWBIKeys() (*WBIKeys, error) {
	// 获取导航信息，其中包含WBI密钥
	navInfo, err := c.GetUserNav()
	if err != nil {
		return nil, fmt.Errorf("failed to get nav info: %w", err)
	}
	
	wbiImg, ok := navInfo["wbi_img"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("wbi_img not found in nav info")
	}
	
	imgURL, ok := wbiImg["img_url"].(string)
	if !ok {
		return nil, fmt.Errorf("img_url not found in wbi_img")
	}
	
	subURL, ok := wbiImg["sub_url"].(string)
	if !ok {
		return nil, fmt.Errorf("sub_url not found in wbi_img")
	}
	
	return &WBIKeys{
		ImgURL: imgURL,
		SubURL: subURL,
	}, nil
}

// GetUserNav 获取用户导航信息 (包含WBI密钥)
func (c *Client) GetUserNav() (map[string]interface{}, error) {
	resp, err := c.httpClient.Get("https://api.bilibili.com/x/web-interface/nav")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 && result.Code != -101 { // -101表示未登录但仍可获取WBI信息
		return nil, fmt.Errorf("get nav info failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	data, ok := result.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid nav data format")
	}
	
	return data, nil
}

// SignParams 为参数添加WBI签名
func (w *WBIManager) SignParams(params map[string]interface{}) (map[string]interface{}, error) {
	// 转换interface{}为string
	stringParams := make(map[string]string)
	for k, v := range params {
		stringParams[k] = fmt.Sprintf("%v", v)
	}
	
	// 使用Sign方法
	signedStringParams := w.Sign(stringParams)
	
	// 转换回interface{}
	result := make(map[string]interface{})
	for k, v := range signedStringParams {
		result[k] = v
	}
	
	return result, nil
}