package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetQRCode 获取登录二维码
func (c *Client) GetQRCode() (*QRCodeResponse, error) {
	timestamp := time.Now().Unix()

	// 构建参数
	params := url.Values{}
	params.Set("appkey", BiliTVAppKey)
	params.Set("local_id", "0")
	params.Set("ts", strconv.FormatInt(timestamp, 10))

	// 签名
	paramStr := params.Encode()
	sign := Sign(paramStr, BiliTVAppSec)
	params.Set("sign", sign)

	// 发送请求
	resp, err := c.httpClient.PostForm("https://passport.bilibili.com/x/passport-tv-login/qrcode/auth_code", params)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	var qrResp QRCodeResponse
	if err := json.Unmarshal(body, &qrResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	return &qrResp, nil
}

// PollQRCode 轮询二维码登录状态
func (c *Client) PollQRCode(authCode string) (*LoginInfo, error) {
	timestamp := time.Now().Unix()

	// 构建参数
	params := url.Values{}
	params.Set("appkey", BiliTVAppKey)
	params.Set("auth_code", authCode)
	params.Set("local_id", "0")
	params.Set("ts", strconv.FormatInt(timestamp, 10))

	// 签名
	paramStr := params.Encode()
	sign := Sign(paramStr, BiliTVAppSec)
	params.Set("sign", sign)

	// 循环轮询
	for {
		resp, err := c.httpClient.PostForm("https://passport.bilibili.com/x/passport-tv-login/qrcode/poll", params)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read response body failed: %w", err)
		}

		var result ResponseData
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("unmarshal response failed: %w", err)
		}

		switch result.Code {
		case 0:
			// 登录成功
			loginData, _ := json.Marshal(result.Data)
			var loginInfo LoginInfo
			if err := json.Unmarshal(loginData, &loginInfo); err != nil {
				return nil, fmt.Errorf("unmarshal login info failed: %w", err)
			}
			loginInfo.Platform = "BiliTV"
			return &loginInfo, nil
		case 86039:
			// 二维码尚未确认，继续轮询
			time.Sleep(1 * time.Second)
			continue
		default:
			return nil, fmt.Errorf("login failed: code=%d, message=%s", result.Code, result.Message)
		}
	}
}



// GetMyInfo 获取当前登录用户的详细信息 (myinfo API)
func (c *Client) GetMyInfo(cookies string) (*MyInfoResponse, error) {
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/space/myinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}

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

	if result.Code != 0 {
		return nil, fmt.Errorf("get my info failed: code=%d, message=%s", result.Code, result.Message)
	}

	// 解析我的信息数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}

	var myInfo MyInfoResponse
	if err := json.Unmarshal(dataBytes, &myInfo); err != nil {
		return nil, fmt.Errorf("unmarshal my info failed: %w", err)
	}

	// 设置兼容性字段
	myInfo.PostProcess()

	return &myInfo, nil
}

// GetMyInfoWithRetry 带重试机制的获取我的信息
func (c *Client) GetMyInfoWithRetry(cookies string, maxRetries int) (*MyInfoResponse, error) {
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var lastErr error
	baseDelay := time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		myInfo, err := c.GetMyInfo(cookies)
		if err == nil {
			return myInfo, nil
		}

		lastErr = err

		// 如果是最后一次尝试，直接返回错误
		if attempt == maxRetries {
			break
		}

		// 检查是否是限流错误
		if IsRateLimitError(err) {
			// 限流错误使用更长的延迟
			delay := time.Duration(attempt*attempt) * baseDelay * 3
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}
			time.Sleep(delay)
		} else {
			// 其他错误使用较短的延迟
			delay := time.Duration(attempt) * baseDelay
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("failed after %d attempts, last error: %w", maxRetries, lastErr)
}

// GetArchivePre 获取投稿预处理信息，包括分区列表
func (c *Client) GetArchivePre(cookies string) (*ArchivePreData, error) {
	req, err := http.NewRequest("GET", "https://member.bilibili.com/x/vupre/web/archive/pre", nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}

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

	if result.Code != 0 {
		return nil, fmt.Errorf("get archive pre failed: code=%d, message=%s", result.Code, result.Message)
	}

	// 解析数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}

	var archiveData ArchivePreData
	if err := json.Unmarshal(dataBytes, &archiveData); err != nil {
		return nil, fmt.Errorf("unmarshal archive data failed: %w", err)
	}

	return &archiveData, nil
}