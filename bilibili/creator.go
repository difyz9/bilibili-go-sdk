package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CreatorInfo 创作者信息
type CreatorInfo struct {
	Mid       int64  `json:"mid"`
	Name      string `json:"name"`
	Face      string `json:"face"`
	Sign      string `json:"sign"`
	Level     int    `json:"level"`
	Silence   int    `json:"silence"`
	VipType   int    `json:"vip_type"`
	VipStatus int    `json:"vip_status"`
}

// DraftInfo 草稿信息
type DraftInfo struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Summary  string `json:"summary"`
	Status   int    `json:"status"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

// TemplateInfo 模板信息
type TemplateInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	Tag         string `json:"tag"`
	Tid         int    `json:"tid"`
	Created     int64  `json:"created"`
}

// GetCreatorInfo 获取创作者信息
func (c *Client) GetCreatorInfo(cookies string) (*CreatorInfo, error) {
	params := make(map[string]interface{})
	
	// 添加WBI签名
	signedParams, err := c.wbiManager.SignParams(params)
	if err != nil {
		return nil, fmt.Errorf("sign params failed: %w", err)
	}
	
	// 构建URL参数
	urlParams := url.Values{}
	for k, v := range signedParams {
		urlParams.Set(k, fmt.Sprintf("%v", v))
	}
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/web/account/info?%s", urlParams.Encode())
	
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
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("get creator info failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析创作者数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var creatorInfo CreatorInfo
	if err := json.Unmarshal(dataBytes, &creatorInfo); err != nil {
		return nil, fmt.Errorf("unmarshal creator info failed: %w", err)
	}
	
	return &creatorInfo, nil
}

// GetDraftList 获取草稿列表
func (c *Client) GetDraftList(page, pageSize int, cookies string) ([]DraftInfo, error) {
	params := make(map[string]interface{})
	params["pn"] = page
	params["ps"] = pageSize
	
	// 添加WBI签名
	signedParams, err := c.wbiManager.SignParams(params)
	if err != nil {
		return nil, fmt.Errorf("sign params failed: %w", err)
	}
	
	// 构建URL参数
	urlParams := url.Values{}
	for k, v := range signedParams {
		urlParams.Set(k, fmt.Sprintf("%v", v))
	}
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/vu/web/draft/list?%s", urlParams.Encode())
	
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
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("get draft list failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析草稿列表数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var drafts []DraftInfo
	if draftArray, ok := dataMap["drafts"].([]interface{}); ok {
		for _, draftItem := range draftArray {
			if draftMap, ok := draftItem.(map[string]interface{}); ok {
				draft := DraftInfo{}
				if id, ok := draftMap["id"].(float64); ok {
					draft.ID = int64(id)
				}
				if title, ok := draftMap["title"].(string); ok {
					draft.Title = title
				}
				if cover, ok := draftMap["cover"].(string); ok {
					draft.Cover = cover
				}
				if summary, ok := draftMap["summary"].(string); ok {
					draft.Summary = summary
				}
				if status, ok := draftMap["status"].(float64); ok {
					draft.Status = int(status)
				}
				if created, ok := draftMap["created"].(float64); ok {
					draft.Created = int64(created)
				}
				if modified, ok := draftMap["modified"].(float64); ok {
					draft.Modified = int64(modified)
				}
				drafts = append(drafts, draft)
			}
		}
	}
	
	return drafts, nil
}

// SaveDraft 保存草稿
func (c *Client) SaveDraft(title, desc, cover, tag string, tid int, cookies string) (int64, error) {
	params := make(map[string]interface{})
	params["title"] = title
	params["desc"] = desc
	params["cover"] = cover
	params["tag"] = tag
	params["tid"] = tid
	params["copyright"] = 1
	params["source"] = ""
	params["dynamic"] = ""
	params["interactive"] = 0
	
	// 添加WBI签名
	signedParams, err := c.wbiManager.SignParams(params)
	if err != nil {
		return 0, fmt.Errorf("sign params failed: %w", err)
	}
	
	// 构建请求参数
	formData := url.Values{}
	for k, v := range signedParams {
		formData.Set(k, fmt.Sprintf("%v", v))
	}
	
	apiURL := "https://member.bilibili.com/x/vu/web/draft/add"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return 0, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return 0, fmt.Errorf("save draft failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析草稿ID
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return 0, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return 0, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var draftID int64
	if id, ok := dataMap["id"].(float64); ok {
		draftID = int64(id)
	}
	
	return draftID, nil
}

// DeleteDraft 删除草稿
func (c *Client) DeleteDraft(draftID int64, cookies string) error {
	params := make(map[string]interface{})
	params["id"] = draftID
	
	// 添加WBI签名
	signedParams, err := c.wbiManager.SignParams(params)
	if err != nil {
		return fmt.Errorf("sign params failed: %w", err)
	}
	
	// 构建请求参数
	formData := url.Values{}
	for k, v := range signedParams {
		formData.Set(k, fmt.Sprintf("%v", v))
	}
	
	apiURL := "https://member.bilibili.com/x/vu/web/draft/delete"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
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
		return fmt.Errorf("delete draft failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	return nil
}

// GetTemplateList 获取模板列表
func (c *Client) GetTemplateList(page, pageSize int, cookies string) ([]TemplateInfo, error) {
	params := make(map[string]interface{})
	params["pn"] = page
	params["ps"] = pageSize
	
	// 添加WBI签名
	signedParams, err := c.wbiManager.SignParams(params)
	if err != nil {
		return nil, fmt.Errorf("sign params failed: %w", err)
	}
	
	// 构建URL参数
	urlParams := url.Values{}
	for k, v := range signedParams {
		urlParams.Set(k, fmt.Sprintf("%v", v))
	}
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/vu/web/template/list?%s", urlParams.Encode())
	
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
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("get template list failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析模板列表数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var templates []TemplateInfo
	if templateArray, ok := dataMap["templates"].([]interface{}); ok {
		for _, templateItem := range templateArray {
			if templateMap, ok := templateItem.(map[string]interface{}); ok {
				template := TemplateInfo{}
				if id, ok := templateMap["id"].(float64); ok {
					template.ID = int64(id)
				}
				if name, ok := templateMap["name"].(string); ok {
					template.Name = name
				}
				if desc, ok := templateMap["description"].(string); ok {
					template.Description = desc
				}
				if cover, ok := templateMap["cover"].(string); ok {
					template.Cover = cover
				}
				if tag, ok := templateMap["tag"].(string); ok {
					template.Tag = tag
				}
				if tid, ok := templateMap["tid"].(float64); ok {
					template.Tid = int(tid)
				}
				if created, ok := templateMap["created"].(float64); ok {
					template.Created = int64(created)
				}
				templates = append(templates, template)
			}
		}
	}
	
	return templates, nil
}

// CreateTemplate 创建模板
func (c *Client) CreateTemplate(name, description, cover, tag string, tid int, cookies string) (int64, error) {
	params := make(map[string]interface{})
	params["name"] = name
	params["description"] = description
	params["cover"] = cover
	params["tag"] = tag
	params["tid"] = tid
	
	// 添加WBI签名
	signedParams, err := c.wbiManager.SignParams(params)
	if err != nil {
		return 0, fmt.Errorf("sign params failed: %w", err)
	}
	
	// 构建请求参数
	formData := url.Values{}
	for k, v := range signedParams {
		formData.Set(k, fmt.Sprintf("%v", v))
	}
	
	apiURL := "https://member.bilibili.com/x/vu/web/template/add"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return 0, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return 0, fmt.Errorf("create template failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析模板ID
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return 0, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return 0, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var templateID int64
	if id, ok := dataMap["id"].(float64); ok {
		templateID = int64(id)
	}
	
	return templateID, nil
}