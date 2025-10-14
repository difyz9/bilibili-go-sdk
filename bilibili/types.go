package bilibili

import (
	"fmt"
	"strings"
)

// ResponseData 通用API响应结构
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TTL     int         `json:"ttl,omitempty"`
}

// QRCodeData 二维码数据
type QRCodeData struct {
	URL      string `json:"url"`
	AuthCode string `json:"auth_code"`
}

// QRCodeResponse 二维码响应
type QRCodeResponse struct {
	Code int        `json:"code"`
	Data QRCodeData `json:"data"`
}

// LoginInfo 登录信息
type LoginInfo struct {
	CookieInfo map[string]interface{} `json:"cookie_info"`
	SSO        []string               `json:"sso"`
	TokenInfo  TokenInfo              `json:"token_info"`
	Platform   string                 `json:"platform,omitempty"`
}

// LoginResponse 登录响应 (兼容短信登录和密码登录)
type LoginResponse struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
	CookieInfo CookieInfo            `json:"cookie_info"`
	TokenInfo  TokenInfo             `json:"token_info"`
}

// CookieInfo Cookie信息
type CookieInfo struct {
	Cookies string `json:"cookies"`
}

// GetCookieString 获取Cookie字符串
func (li *LoginInfo) GetCookieString() string {
	cookies, ok := li.CookieInfo["cookies"].([]interface{})
	if !ok {
		return ""
	}

	var cookieStrs []string
	for _, cookie := range cookies {
		if cookieMap, ok := cookie.(map[string]interface{}); ok {
			if name, nameOk := cookieMap["name"].(string); nameOk {
				if value, valueOk := cookieMap["value"].(string); valueOk {
					cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", name, value))
				}
			}
		}
	}

	return strings.Join(cookieStrs, "; ")
}

// GetCSRFToken 获取CSRF token
func (li *LoginInfo) GetCSRFToken() (string, error) {
	cookies, ok := li.CookieInfo["cookies"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no cookies found")
	}

	for _, cookie := range cookies {
		if cookieMap, ok := cookie.(map[string]interface{}); ok {
			if name, nameOk := cookieMap["name"].(string); nameOk && name == "bili_jct" {
				if value, valueOk := cookieMap["value"].(string); valueOk {
					return value, nil
				}
			}
		}
	}

	return "", fmt.Errorf("CSRF token not found in cookie")
}

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Mid          int64  `json:"mid"`
	RefreshToken string `json:"refresh_token"`
	Uname        string `json:"uname,omitempty"`
	Face         string `json:"face,omitempty"`
}

// UserBasicInfo 用户基本信息
type UserBasicInfo struct {
	Mid      int64  `json:"mid"`
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Face     string `json:"face"`
	Sign     string `json:"sign"`
	Rank     int    `json:"rank"`
	Level    int    `json:"level"`
	Birthday string `json:"birthday"`
}

// MyInfoResponse 详细用户信息响应结构 (myinfo API)
type MyInfoResponse struct {
	Mid       int64  `json:"mid"`
	Uname     string `json:"uname"`
	UserID    string `json:"userid"`
	Sign      string `json:"sign"`
	Birthday  string `json:"birthday"`
	Sex       string `json:"sex"`
	NickFree  bool   `json:"nick_free"`
	Rank      string `json:"rank"`
	Face      string `json:"face"`
	Level     int    `json:"level"`
	Silence   int    `json:"silence"`
	Coins     int    `json:"coins"`
	Fans      int    `json:"fans"`
	Friend    int    `json:"friend"`
	Attention int    `json:"attention"`
}

// PartitionType 分区类型
type PartitionType struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Desc     string          `json:"desc,omitempty"`
	Children []PartitionType `json:"children,omitempty"`
}

// ArchivePreData archive_pre接口返回的数据
type ArchivePreData struct {
	TypeList []PartitionType `json:"typelist"`
}