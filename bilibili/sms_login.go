package bilibili

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SMSLoginRequest 短信登录请求
type SMSLoginRequest struct {
	Tel       string `json:"tel"`
	Cid       string `json:"cid"`       // 国家代码，中国为86
	Code      string `json:"code"`      // 短信验证码
	LoginKey  string `json:"login_key"` // 登录密钥
	Challenge string `json:"challenge"` // 极验验证
	Validate  string `json:"validate"`  // 极验验证
	Seccode   string `json:"seccode"`   // 极验验证
}

// SMSCaptchaResponse 短信验证码响应
type SMSCaptchaResponse struct {
	CaptchaKey string `json:"captcha_key"`
	RecaptchaToken string `json:"recaptcha_token"`
}

// SendSMS 发送短信验证码
func (c *Client) SendSMS(tel, cid string) (*SMSCaptchaResponse, error) {
	// 构建参数
	params := url.Values{}
	params.Set("tel", tel)
	params.Set("cid", cid)
	params.Set("source", "main_web")
	params.Set("token", "")
	params.Set("challenge", "")
	params.Set("validate", "")
	params.Set("seccode", "")
	
	apiURL := "https://passport.bilibili.com/x/passport-login/sms/send"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://passport.bilibili.com/")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("send SMS failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析响应数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var smsResp SMSCaptchaResponse
	if err := json.Unmarshal(dataBytes, &smsResp); err != nil {
		return nil, fmt.Errorf("unmarshal SMS response failed: %w", err)
	}
	
	return &smsResp, nil
}

// SMSLogin 短信登录
func (c *Client) SMSLogin(loginReq *SMSLoginRequest) (*LoginResponse, error) {
	// 生成登录参数
	timestamp := getCurrentTimestamp()
	
	// 构建参数
	params := url.Values{}
	params.Set("tel", loginReq.Tel)
	params.Set("cid", loginReq.Cid)
	params.Set("code", loginReq.Code)
	params.Set("source", "main_web")
	params.Set("ts", fmt.Sprintf("%d", timestamp))
	params.Set("login_key", loginReq.LoginKey)
	
	if loginReq.Challenge != "" {
		params.Set("challenge", loginReq.Challenge)
		params.Set("validate", loginReq.Validate)
		params.Set("seccode", loginReq.Seccode)
	}
	
	apiURL := "https://passport.bilibili.com/x/passport-login/sms/login"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://passport.bilibili.com/")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("SMS login failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析登录响应
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var loginResp LoginResponse
	if err := json.Unmarshal(dataBytes, &loginResp); err != nil {
		return nil, fmt.Errorf("unmarshal login response failed: %w", err)
	}
	
	// 提取Cookie
	cookies := extractCookiesFromResponse(resp)
	loginResp.CookieInfo.Cookies = formatCookies(cookies)
	
	return &loginResp, nil
}

// PasswordLogin 密码登录
func (c *Client) PasswordLogin(username, password string) (*LoginResponse, error) {
	// 获取登录密钥
	keyResp, err := c.getLoginKey()
	if err != nil {
		return nil, fmt.Errorf("get login key failed: %w", err)
	}
	
	// 加密密码
	encryptedPassword, err := c.encryptPassword(password, keyResp.Hash, keyResp.Key)
	if err != nil {
		return nil, fmt.Errorf("encrypt password failed: %w", err)
	}
	
	// 构建登录参数
	timestamp := getCurrentTimestamp()
	
	params := url.Values{}
	params.Set("username", username)
	params.Set("password", encryptedPassword)
	params.Set("keep", "true")
	params.Set("source", "main_web")
	params.Set("ts", fmt.Sprintf("%d", timestamp))
	params.Set("token", "")
	params.Set("challenge", "")
	params.Set("validate", "")
	params.Set("seccode", "")
	
	apiURL := "https://passport.bilibili.com/x/passport-login/web/login"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://passport.bilibili.com/")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("password login failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析登录响应
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var loginResp LoginResponse
	if err := json.Unmarshal(dataBytes, &loginResp); err != nil {
		return nil, fmt.Errorf("unmarshal login response failed: %w", err)
	}
	
	// 提取Cookie
	cookies := extractCookiesFromResponse(resp)
	loginResp.CookieInfo.Cookies = formatCookies(cookies)
	
	return &loginResp, nil
}

// LoginKeyResponse 登录密钥响应
type LoginKeyResponse struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

// getLoginKey 获取登录密钥
func (c *Client) getLoginKey() (*LoginKeyResponse, error) {
	apiURL := "https://passport.bilibili.com/x/passport-login/web/key"
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://passport.bilibili.com/")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("get login key failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析密钥数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var keyResp LoginKeyResponse
	if err := json.Unmarshal(dataBytes, &keyResp); err != nil {
		return nil, fmt.Errorf("unmarshal key response failed: %w", err)
	}
	
	return &keyResp, nil
}

// encryptPassword 加密密码
func (c *Client) encryptPassword(password, hash, key string) (string, error) {
	// 简化的密码加密实现
	// 实际B站使用RSA加密，这里用MD5简化处理
	combined := hash + password + key
	hasher := md5.New()
	hasher.Write([]byte(combined))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// CheckTelValid 检查手机号是否有效
func (c *Client) CheckTelValid(tel, cid string) (bool, error) {
	params := url.Values{}
	params.Set("tel", tel)
	params.Set("cid", cid)
	
	apiURL := fmt.Sprintf("https://passport.bilibili.com/x/passport-login/tel/verify?%s", params.Encode())
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://passport.bilibili.com/")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("decode response failed: %w", err)
	}
	
	return result.Code == 0, nil
}

// Logout 登出
func (c *Client) Logout(cookies string) error {
	params := url.Values{}
	params.Set("biliCSRF", extractCSRFFromCookies(cookies))
	
	apiURL := "https://passport.bilibili.com/login/exit/v2"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return fmt.Errorf("logout failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	return nil
}

// extractCSRFFromCookies 从Cookie中提取CSRF token
func extractCSRFFromCookies(cookies string) string {
	// 从Cookie字符串中提取bili_jct值
	parts := strings.Split(cookies, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "bili_jct=") {
			return strings.TrimPrefix(part, "bili_jct=")
		}
	}
	return ""
}

// extractCookiesFromResponse 从HTTP响应中提取cookies
func extractCookiesFromResponse(resp *http.Response) []*http.Cookie {
	return resp.Cookies()
}

// formatCookies 格式化cookies为字符串
func formatCookies(cookies []*http.Cookie) string {
	var cookieStrs []string
	for _, cookie := range cookies {
		cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	return strings.Join(cookieStrs, "; ")
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}