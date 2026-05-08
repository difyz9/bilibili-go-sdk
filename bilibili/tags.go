package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TagInfo 标签信息
type TagInfo struct {
	Tag         string `json:"tag,omitempty"`
	Name        string `json:"name"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	State       int    `json:"state"`
	Checked     int    `json:"checked,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
}

type tagCheckResponse struct {
	Code    int `json:"code"`
	Data    struct {
		Code    int    `json:"code"`
		Content string `json:"content"`
	} `json:"data"`
	Message string `json:"message"`
}

// TagRecommendRequest 标签推荐请求
type TagRecommendRequest struct {
	UploadID    string
	SubtypeID   int
	Title       string
	Filename    string
	Description string
	CoverURL    string
}

// CheckTag 检查标签是否有效
func (c *Client) CheckTag(tag string) (bool, error) {
	escapedTag := url.QueryEscape(tag)
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/vupre/web/topic/tag/check?tag=%s", escapedTag)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return false, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result tagCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return false, fmt.Errorf("check tag failed: code=%d, message=%s", result.Code, result.Message)
	}

	return result.Data.Code == 0, nil
}

// GetRecommendedTags 获取推荐标签
func (c *Client) GetRecommendedTags(title, description string, cookies string) ([]TagInfo, error) {
	return c.RecommendTags(&TagRecommendRequest{
		Title:       title,
		Filename:    title,
		Description: description,
	}, cookies)
}

// RecommendTags 使用创作中心接口获取推荐标签
func (c *Client) RecommendTags(request *TagRecommendRequest, cookies string) ([]TagInfo, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	params := url.Values{}
	if request.UploadID != "" {
		params.Set("upload_id", request.UploadID)
	}
	if request.SubtypeID > 0 {
		params.Set("subtype_id", fmt.Sprintf("%d", request.SubtypeID))
	}
	if request.Title != "" {
		params.Set("title", request.Title)
	}
	if request.Filename != "" {
		params.Set("filename", request.Filename)
	}
	if request.Description != "" {
		params.Set("description", request.Description)
	}
	if request.CoverURL != "" {
		params.Set("cover_url", strings.TrimPrefix(strings.TrimPrefix(request.CoverURL, "https://"), "http://"))
	}
	params.Set("t", fmt.Sprintf("%d", time.Now().UnixMilli()))
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/vupre/web/tag/recommend?%s", params.Encode())
	
	req, err := http.NewRequest("GET", apiURL, nil)
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
	
	var result struct {
		Code      int       `json:"code"`
		Data      []TagInfo `json:"data"`
		Message   string    `json:"message"`
		RequestID string    `json:"request_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("get recommended tags failed: code=%d, message=%s", result.Code, result.Message)
	}

	for index := range result.Data {
		if result.Data[index].Name == "" {
			result.Data[index].Name = result.Data[index].Tag
		}
		result.Data[index].RequestID = result.RequestID
	}
	
	return result.Data, nil
}

// SearchTags 搜索标签
func (c *Client) SearchTags(keyword string, cookies string) ([]TagInfo, error) {
	// 这里可以实现基于关键词的标签搜索
	// 目前B站没有公开的标签搜索API，可以基于推荐标签功能实现
	return c.GetRecommendedTags(keyword, keyword, cookies)
}

// ValidateTags 批量验证标签
func (c *Client) ValidateTags(tags []string) (map[string]bool, error) {
	result := make(map[string]bool)
	
	for _, tag := range tags {
		valid, err := c.CheckTag(tag)
		if err != nil {
			return nil, fmt.Errorf("failed to check tag %s: %w", tag, err)
		}
		result[tag] = valid
	}
	
	return result, nil
}

// FormatTags 格式化标签为B站要求的格式
func FormatTags(tags []string) string {
	// 过滤空标签和过长标签
	var validTags []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" && len(tag) <= 20 { // B站标签长度限制
			validTags = append(validTags, tag)
		}
	}
	
	// 限制标签数量 (B站限制10个标签)
	if len(validTags) > 10 {
		validTags = validTags[:10]
	}
	
	return strings.Join(validTags, ",")
}