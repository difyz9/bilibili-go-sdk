package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetUserBasicInfo 获取指定用户的空间详细信息。
//
// 接口: GET https://api.bilibili.com/x/space/wbi/acc/info
// 需要提供包含 SESSDATA 的 Cookie，且请求参数需要附带 WBI 签名。
func (c *Client) GetUserBasicInfo(mid int64, cookies string) (*UserBasicInfo, error) {
	if mid <= 0 {
		return nil, fmt.Errorf("mid must be greater than 0")
	}

	signedParams, err := c.SignWithWBI(map[string]string{
		"mid": strconv.FormatInt(mid, 10),
	})
	if err != nil {
		return nil, fmt.Errorf("sign params failed: %w", err)
	}

	query := url.Values{}
	for key, value := range signedParams {
		query.Set(key, value)
	}

	apiURL := fmt.Sprintf("https://api.bilibili.com/x/space/wbi/acc/info?%s", query.Encode())
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://space.bilibili.com/")
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}

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
		return nil, fmt.Errorf("get user basic info failed: code=%d, message=%s", result.Code, result.Message)
	}

	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}

	var userInfo UserBasicInfo
	if err := json.Unmarshal(dataBytes, &userInfo); err != nil {
		return nil, fmt.Errorf("unmarshal user basic info failed: %w", err)
	}

	return &userInfo, nil
}
