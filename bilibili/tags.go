package bilibili

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// TagInfo 标签信息
type TagInfo struct {
	Name        string `json:"name"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	State       int    `json:"state"`
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
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("decode response failed: %w", err)
	}
	
	return result.Code == 0, nil
}

// GetRecommendedTags 获取推荐标签
func (c *Client) GetRecommendedTags(title, description string, cookies string) ([]TagInfo, error) {
	params := url.Values{}
	params.Set("title", title)
	params.Set("filename", title)
	params.Set("desc", description)
	params.Set("cover", "")
	params.Set("groupid", "1")
	params.Set("vfea", "")
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/web/archive/tags?%s", params.Encode())
	
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
	
	var result ResponseData
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("get recommended tags failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析标签数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var tagsData map[string]interface{}
	if err := json.Unmarshal(dataBytes, &tagsData); err != nil {
		return nil, fmt.Errorf("unmarshal tags data failed: %w", err)
	}
	
	var tags []TagInfo
	if tagsArray, ok := tagsData["tags"].([]interface{}); ok {
		for _, tagItem := range tagsArray {
			if tagMap, ok := tagItem.(map[string]interface{}); ok {
				tag := TagInfo{}
				if name, ok := tagMap["name"].(string); ok {
					tag.Name = name
				}
				if cover, ok := tagMap["cover"].(string); ok {
					tag.Cover = cover
				}
				if desc, ok := tagMap["description"].(string); ok {
					tag.Description = desc
				}
				if tagType, ok := tagMap["type"].(float64); ok {
					tag.Type = int(tagType)
				}
				if state, ok := tagMap["state"].(float64); ok {
					tag.State = int(state)
				}
				tags = append(tags, tag)
			}
		}
	}
	
	return tags, nil
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