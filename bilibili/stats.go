package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// VideoStat 视频统计信息
type VideoStat struct {
	AID       int64 `json:"aid"`
	BVid      string `json:"bvid"`
	View      int   `json:"view"`
	Danmaku   int   `json:"danmaku"`
	Reply     int   `json:"reply"`
	Favorite  int   `json:"favorite"`
	Coin      int   `json:"coin"`
	Share     int   `json:"share"`
	Like      int   `json:"like"`
	NowRank   int   `json:"now_rank"`
	HisRank   int   `json:"his_rank"`
	NoReprint int   `json:"no_reprint"`
	Copyright int   `json:"copyright"`
}

// UserStat 用户统计信息
type UserStat struct {
	Mid       int64 `json:"mid"`
	Following int   `json:"following"`
	Whisper   int   `json:"whisper"`
	Black     int   `json:"black"`
	Follower  int   `json:"follower"`
}

// UpStat UP主统计信息
type UpStat struct {
	Archive ArchiveStat `json:"archive"`
	Article ArticleStat `json:"article"`
	Likes   int         `json:"likes"`
}

// ArchiveStat 投稿统计
type ArchiveStat struct {
	View int `json:"view"`
}

// ArticleStat 专栏统计
type ArticleStat struct {
	View int `json:"view"`
}

// GetVideoStat 获取视频统计信息
func (c *Client) GetVideoStat(bvid string) (*VideoStat, error) {
	apiURL := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvid)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	
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
		return nil, fmt.Errorf("get video stat failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析视频数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var videoStat VideoStat
	if aid, ok := dataMap["aid"].(float64); ok {
		videoStat.AID = int64(aid)
	}
	if bvid, ok := dataMap["bvid"].(string); ok {
		videoStat.BVid = bvid
	}
	if stat, ok := dataMap["stat"].(map[string]interface{}); ok {
		if view, ok := stat["view"].(float64); ok {
			videoStat.View = int(view)
		}
		if danmaku, ok := stat["danmaku"].(float64); ok {
			videoStat.Danmaku = int(danmaku)
		}
		if reply, ok := stat["reply"].(float64); ok {
			videoStat.Reply = int(reply)
		}
		if favorite, ok := stat["favorite"].(float64); ok {
			videoStat.Favorite = int(favorite)
		}
		if coin, ok := stat["coin"].(float64); ok {
			videoStat.Coin = int(coin)
		}
		if share, ok := stat["share"].(float64); ok {
			videoStat.Share = int(share)
		}
		if like, ok := stat["like"].(float64); ok {
			videoStat.Like = int(like)
		}
	}
	
	return &videoStat, nil
}

// GetUserStat 获取用户统计信息
func (c *Client) GetUserStat(mid int64) (*UserStat, error) {
	apiURL := fmt.Sprintf("https://api.bilibili.com/x/relation/stat?vmid=%d", mid)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	
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
		return nil, fmt.Errorf("get user stat failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析用户统计数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var userStat UserStat
	if err := json.Unmarshal(dataBytes, &userStat); err != nil {
		return nil, fmt.Errorf("unmarshal user stat failed: %w", err)
	}
	
	userStat.Mid = mid
	return &userStat, nil
}

// GetUpStat 获取UP主统计信息
func (c *Client) GetUpStat(mid int64) (*UpStat, error) {
	apiURL := fmt.Sprintf("https://api.bilibili.com/x/space/upstat?mid=%d", mid)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	
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
		return nil, fmt.Errorf("get up stat failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析UP主统计数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var upStat UpStat
	if err := json.Unmarshal(dataBytes, &upStat); err != nil {
		return nil, fmt.Errorf("unmarshal up stat failed: %w", err)
	}
	
	return &upStat, nil
}

// GetMyUpStat 获取我的UP主统计信息（需要登录）
func (c *Client) GetMyUpStat(cookies string) (*UpStat, error) {
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
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/space/myinfo?%s", urlParams.Encode())
	
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
		return nil, fmt.Errorf("get my up stat failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析UP主统计数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var dataMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &dataMap); err != nil {
		return nil, fmt.Errorf("unmarshal data failed: %w", err)
	}
	
	var upStat UpStat
	if stat, ok := dataMap["stat"].(map[string]interface{}); ok {
		statBytes, _ := json.Marshal(stat)
		json.Unmarshal(statBytes, &upStat)
	}
	
	return &upStat, nil
}

// GetVideoAnalytics 获取视频详细分析数据（需要登录且为视频作者）
func (c *Client) GetVideoAnalytics(aid int64, cookies string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	params["aid"] = aid
	params["period"] = "all" // 全部时间
	
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
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/h5/data/arc/stat?%s", urlParams.Encode())
	
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
		return nil, fmt.Errorf("get video analytics failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析分析数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var analytics map[string]interface{}
	if err := json.Unmarshal(dataBytes, &analytics); err != nil {
		return nil, fmt.Errorf("unmarshal analytics failed: %w", err)
	}
	
	return analytics, nil
}

// GetTrendData 获取趋势数据
func (c *Client) GetTrendData(aid int64, startDate, endDate time.Time, cookies string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	params["aid"] = aid
	params["s_date"] = startDate.Format("2006-01-02")
	params["e_date"] = endDate.Format("2006-01-02")
	
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
	
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/h5/data/trend?%s", urlParams.Encode())
	
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
		return nil, fmt.Errorf("get trend data failed: code=%d, message=%s", result.Code, result.Message)
	}
	
	// 解析趋势数据
	dataBytes, err := json.Marshal(result.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data failed: %w", err)
	}
	
	var trendData map[string]interface{}
	if err := json.Unmarshal(dataBytes, &trendData); err != nil {
		return nil, fmt.Errorf("unmarshal trend data failed: %w", err)
	}
	
	return trendData, nil
}