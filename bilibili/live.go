package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// LiveRoomInfo 直播间信息
type LiveRoomInfo struct {
	RoomID       int64  `json:"room_id"`
	ShortID      int    `json:"short_id"`
	UID          int64  `json:"uid"`
	NeedP2P      int    `json:"need_p2p"`
	IsHidden     bool   `json:"is_hidden"`
	IsLocked     bool   `json:"is_locked"`
	IsPortrait   bool   `json:"is_portrait"`
	LiveStatus   int    `json:"live_status"`
	HiddenTill   int64  `json:"hidden_till"`
	LockTill     int64  `json:"lock_till"`
	Encrypted    bool   `json:"encrypted"`
	PwdVerified  bool   `json:"pwd_verified"`
	LiveTime     int64  `json:"live_time"`
	RoomShield   int    `json:"room_shield"`
	AllSpecialType []int `json:"all_special_type"`
}

// LiveStreamInfo 直播流信息
type LiveStreamInfo struct {
	CurrentQuality    int                    `json:"current_quality"`
	AcceptQuality     []int                 `json:"accept_quality"`
	CurrentQualityName string               `json:"current_qn"`
	QualityDescription []QualityDescription `json:"quality_description"`
	Durl              []StreamURL          `json:"durl"`
}

// QualityDescription 画质描述
type QualityDescription struct {
	Quality int    `json:"qn"`
	Desc    string `json:"desc"`
	HDDesc  string `json:"hd_desc"`
}

// StreamURL 流地址
type StreamURL struct {
	URL       string `json:"url"`
	Length    int64  `json:"length"`
	Order     int    `json:"order"`
	StreamType int   `json:"stream_type"`
	P2PType   int    `json:"p2p_type"`
}

// GetLiveRoomInfo 获取直播间信息
func (c *Client) GetLiveRoomInfo(roomID int64) (*LiveRoomInfo, error) {
	apiURL := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/get_info?room_id=%d", roomID)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://live.bilibili.com/")
	
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
		return nil, fmt.Errorf("get live room info failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析直播间数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var roomInfo LiveRoomInfo
	if err := json.Unmarshal(dataBytes, &roomInfo); err != nil {
		return nil, fmt.Errorf("unmarshal room info failed: %w", err)
	}
	
	return &roomInfo, nil
}

// GetLiveStreamInfo 获取直播流信息
func (c *Client) GetLiveStreamInfo(roomID int64, quality int, cookies string) (*LiveStreamInfo, error) {
	params := url.Values{}
	params.Set("cid", fmt.Sprintf("%d", roomID))
	params.Set("qn", fmt.Sprintf("%d", quality))
	params.Set("platform", "web")
	params.Set("https_url_req", "1")
	params.Set("ptype", "8")
	
	apiURL := fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo?%s", params.Encode())
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://live.bilibili.com/")
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
		return nil, fmt.Errorf("get live stream info failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析流数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var streamInfo LiveStreamInfo
	if playurlInfo, ok := dataMap["playurl_info"].(map[string]interface{}); ok {
		if playurl, ok := playurlInfo["playurl"].(map[string]interface{}); ok {
			playurlBytes, _ := json.Marshal(playurl)
			json.Unmarshal(playurlBytes, &streamInfo)
		}
	}
	
	return &streamInfo, nil
}

// StartLive 开始直播
func (c *Client) StartLive(roomID int64, area int, cookies string) error {
	params := make(map[string]interface{})
	params["room_id"] = roomID
	params["area_v2"] = area
	params["platform"] = "pc"
	
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
	
	apiURL := "https://api.live.bilibili.com/room/v1/Room/startLive"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://live.bilibili.com/")
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
		return fmt.Errorf("start live failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	return nil
}

// StopLive 停止直播
func (c *Client) StopLive(roomID int64, cookies string) error {
	params := make(map[string]interface{})
	params["room_id"] = roomID
	
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
	
	apiURL := "https://api.live.bilibili.com/room/v1/Room/stopLive"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://live.bilibili.com/")
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
		return fmt.Errorf("stop live failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	return nil
}

// UpdateLiveTitle 更新直播标题
func (c *Client) UpdateLiveTitle(roomID int64, title string, cookies string) error {
	params := make(map[string]interface{})
	params["room_id"] = roomID
	params["title"] = title
	
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
	
	apiURL := "https://api.live.bilibili.com/room/v1/Room/update"
	
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://live.bilibili.com/")
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
		return fmt.Errorf("update live title failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	return nil
}