package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// VideoEditRequest 视频编辑请求
type VideoEditRequest struct {
	AID         int64  `json:"aid"`
	BVid        string `json:"bvid"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Cover       string `json:"cover"`
	Tag         string `json:"tag"`
	Tid         int    `json:"tid"`
	Copyright   int    `json:"copyright"`
	Source      string `json:"source"`
	OpenElec    int    `json:"open_elec"`
	NoReprint   int    `json:"no_reprint"`
	SelectiOn   int    `json:"selection"`
	Dynamic     string `json:"dynamic"`
	Interactive int    `json:"interactive"`
}

// VideoInfo 视频信息
type VideoInfo struct {
	AID       int64  `json:"aid"`
	BVid      string `json:"bvid"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Cover     string `json:"pic"`
	Tag       string `json:"tag"`
	Tid       int    `json:"tid"`
	Copyright int    `json:"copyright"`
	Source    string `json:"source"`
	State     int    `json:"state"`
	Duration  int    `json:"duration"`
	PubDate   int64  `json:"pubdate"`
	CTime     int64  `json:"ctime"`
}

// GetVideoInfo 获取视频信息
func (c *Client) GetVideoInfo(bvid string, cookies string) (*VideoInfo, error) {
	apiURL := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvid)
	
	req, err := http.NewRequest("GET", apiURL, nil)
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
	
	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}
	
	if result.Code != 0 {
		return nil, fmt.Errorf("get video info failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析视频数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var videoInfo VideoInfo
	if err := json.Unmarshal(dataBytes, &videoInfo); err != nil {
		return nil, fmt.Errorf("unmarshal video info failed: %w", err)
	}
	
	return &videoInfo, nil
}

// EditVideo 编辑视频信息
func (c *Client) EditVideo(editReq *VideoEditRequest, cookies string) error {
	// 使用WBI签名
	params := make(map[string]interface{})
	params["aid"] = editReq.AID
	if editReq.BVid != "" {
		params["bvid"] = editReq.BVid
	}
	params["title"] = editReq.Title
	params["desc"] = editReq.Desc
	params["cover"] = editReq.Cover
	params["tag"] = editReq.Tag
	params["tid"] = editReq.Tid
	params["copyright"] = editReq.Copyright
	params["source"] = editReq.Source
	params["open_elec"] = editReq.OpenElec
	params["no_reprint"] = editReq.NoReprint
	params["selection"] = editReq.SelectiOn
	params["dynamic"] = editReq.Dynamic
	params["interactive"] = editReq.Interactive
	
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
	
	apiURL := "https://member.bilibili.com/x/vu/web/edit"
	
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
		return fmt.Errorf("edit video failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	return nil
}

// DeleteVideo 删除视频
func (c *Client) DeleteVideo(aid int64, cookies string) error {
	params := make(map[string]interface{})
	params["aid"] = aid
	
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
	
	apiURL := "https://member.bilibili.com/x/vu/web/delete"
	
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
		return fmt.Errorf("delete video failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	return nil
}

// GetMyVideos 获取我的视频列表
func (c *Client) GetMyVideos(page, pageSize int, cookies string) ([]VideoInfo, error) {
	params := make(map[string]interface{})
	params["pn"] = page
	params["ps"] = pageSize
	params["order"] = "pubdate"
	params["tid"] = 0
	params["keyword"] = ""
	
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
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/space/mylist?%s", urlParams.Encode())
	
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
		return nil, fmt.Errorf("get my videos failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析视频列表数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var videos []VideoInfo
	if vlistArray, ok := dataMap["vlist"].([]interface{}); ok {
		for _, videoItem := range vlistArray {
			if videoMap, ok := videoItem.(map[string]interface{}); ok {
				video := VideoInfo{}
				if aid, ok := videoMap["aid"].(float64); ok {
					video.AID = int64(aid)
				}
				if bvid, ok := videoMap["bvid"].(string); ok {
					video.BVid = bvid
				}
				if title, ok := videoMap["title"].(string); ok {
					video.Title = title
				}
				if desc, ok := videoMap["description"].(string); ok {
					video.Desc = desc
				}
				if pic, ok := videoMap["pic"].(string); ok {
					video.Cover = pic
				}
				if state, ok := videoMap["state"].(float64); ok {
					video.State = int(state)
				}
				if duration, ok := videoMap["duration"].(float64); ok {
					video.Duration = int(duration)
				}
				if created, ok := videoMap["created"].(float64); ok {
					video.CTime = int64(created)
				}
				videos = append(videos, video)
			}
		}
	}
	
	return videos, nil
}