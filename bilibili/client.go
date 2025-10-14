package bilibili

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// API常量
const (
	BiliTVAppKey  = "4409e2ce8ffd12b8"
	BiliTVAppSec  = "59b43e04ad6965f34319062b478f83dd"
	AndroidAppKey = "783bbb7264451d82"
	AndroidAppSec = "2653583c8873dea268ab9386918b1d65"
)

// Client Bilibili API 客户端
type Client struct {
	httpClient *http.Client
	userAgent  string
	config     *Config
	wbiManager *WBIManager
}

// NewClient 创建新的Bilibili客户端
func NewClient(opts ...Option) *Client {
	config := DefaultConfig()
	config.ApplyOptions(opts...)

	return &Client{
		httpClient: config.HTTPClient,
		userAgent:  config.UserAgent,
		config:     config,
		wbiManager: NewWBIManager(),
	}
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// GetHTTPClient 获取内部HTTP客户端（用于高级配置）
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// SetUserAgent 设置User-Agent
func (c *Client) SetUserAgent(ua string) {
	c.userAgent = ua
}

// Sign 计算API签名
func Sign(params string, appSec string) string {
	h := md5.New()
	h.Write([]byte(params + appSec))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// IsRateLimitError 判断是否是限流错误
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := err.Error()
	return strings.Contains(errMsg, "code=-799") ||
		   strings.Contains(errMsg, "请求过于频繁") ||
		   strings.Contains(errMsg, "rate limit") ||
		   strings.Contains(errMsg, "too many requests")
}

// IsNetworkError 判断是否是网络错误
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}
	
	errorStr := strings.ToLower(err.Error())
	return strings.Contains(errorStr, "broken pipe") ||
		   strings.Contains(errorStr, "connection reset") ||
		   strings.Contains(errorStr, "timeout") ||
		   strings.Contains(errorStr, "network") ||
		   strings.Contains(errorStr, "dial") ||
		   strings.Contains(errorStr, "eof")
}

// GetWBIManager 获取WBI管理器 (用于手动管理WBI签名)
func (c *Client) GetWBIManager() *WBIManager {
	return c.wbiManager
}

// UpdateWBIKeys 更新WBI密钥 (自动获取最新密钥)
func (c *Client) UpdateWBIKeys() error {
	keys, err := c.GetWBIKeys()
	if err != nil {
		return fmt.Errorf("failed to get WBI keys: %w", err)
	}
	
	return c.wbiManager.UpdateKeys(keys.ImgURL, keys.SubURL)
}

// SignWithWBI 使用WBI签名参数 (自动更新过期密钥)
func (c *Client) SignWithWBI(params map[string]string) (map[string]string, error) {
	// 检查密钥是否过期，如果过期则自动更新
	if c.wbiManager.IsExpired() {
		if err := c.UpdateWBIKeys(); err != nil {
			// 更新失败时返回未签名的参数
			return params, fmt.Errorf("failed to update WBI keys: %w", err)
		}
	}
	
	return c.wbiManager.Sign(params), nil
}