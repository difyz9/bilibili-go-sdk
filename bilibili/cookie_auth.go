package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// LoginByCookies 通过Cookie字符串进行登录验证
func (c *Client) LoginByCookies(cookies string) (*LoginInfo, error) {
	// 验证Cookie有效性
	userInfo, err := c.ValidateCookies(cookies)
	if err != nil {
		return nil, fmt.Errorf("cookie validation failed: %w", err)
	}
	
	// 构造LoginInfo
	loginInfo := &LoginInfo{
		CookieInfo: map[string]interface{}{
			"cookies": parseCookieString(cookies),
		},
		TokenInfo: TokenInfo{
			Mid:   userInfo.Mid,
			Uname: userInfo.Name,
			Face:  userInfo.Face,
		},
		Platform: "Cookie",
	}
	
	return loginInfo, nil
}

// ValidateCookies 验证Cookie有效性并返回用户信息
func (c *Client) ValidateCookies(cookies string) (*UserBasicInfo, error) {
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/nav", nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Cookie", cookies)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	
	var result ResponseData
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	if result.Code == -101 {
		return nil, fmt.Errorf("cookies invalid or expired")
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("validate cookies failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析用户信息
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var navData map[string]interface{}
	if err := json.Unmarshal(dataBytes, &navData); err != nil {
		return nil, fmt.Errorf("unmarshal nav data failed: %w", err)
	}
	
	// 检查登录状态
	isLogin, ok := navData["isLogin"].(bool)
	if !ok || !isLogin {
		return nil, fmt.Errorf("user not logged in")
	}
	
	// 提取用户信息
	userInfo := &UserBasicInfo{}
	if mid, ok := navData["mid"].(float64); ok {
		userInfo.Mid = int64(mid)
	}
	if uname, ok := navData["uname"].(string); ok {
		userInfo.Name = uname
	}
	if face, ok := navData["face"].(string); ok {
		userInfo.Face = face
	}
	if level, ok := navData["level_info"].(map[string]interface{}); ok {
		if currentLevel, ok := level["current_level"].(float64); ok {
			userInfo.Level = int(currentLevel)
		}
	}
	
	return userInfo, nil
}

// parseCookieString 解析Cookie字符串为结构化格式
func parseCookieString(cookieStr string) []map[string]interface{} {
	var cookies []map[string]interface{}
	
	pairs := strings.Split(cookieStr, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		cookie := map[string]interface{}{
			"name":  strings.TrimSpace(parts[0]),
			"value": strings.TrimSpace(parts[1]),
		}
		cookies = append(cookies, cookie)
	}
	
	return cookies
}

// LoadCookiesFromFile 从文件加载Cookie (支持多种格式)
func (c *Client) LoadCookiesFromFile(filePath string) (*LoginInfo, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read cookie file failed: %w", err)
	}
	
	// 尝试解析为JSON格式 (biliup格式)
	var cookieData map[string]interface{}
	if err := json.Unmarshal(data, &cookieData); err == nil {
		if cookieInfo, ok := cookieData["cookie_info"]; ok {
			if cookies, ok := cookieInfo.(map[string]interface{})["cookies"]; ok {
				if cookieList, ok := cookies.([]interface{}); ok {
					var cookieStr strings.Builder
					for _, cookie := range cookieList {
						if cookieMap, ok := cookie.(map[string]interface{}); ok {
							if name, nameOk := cookieMap["name"].(string); nameOk {
								if value, valueOk := cookieMap["value"].(string); valueOk {
									if cookieStr.Len() > 0 {
										cookieStr.WriteString("; ")
									}
									cookieStr.WriteString(fmt.Sprintf("%s=%s", name, value))
								}
							}
						}
					}
					return c.LoginByCookies(cookieStr.String())
				}
			}
		}
	}
	
	// 尝试解析为纯文本Cookie字符串
	cookieStr := strings.TrimSpace(string(data))
	if cookieStr != "" {
		return c.LoginByCookies(cookieStr)
	}
	
	return nil, fmt.Errorf("unsupported cookie file format")
}