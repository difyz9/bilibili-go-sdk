package bilibili

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type creativeCenterResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	TTL     int    `json:"ttl,omitempty"`
}

// SeasonListParams 获取合集列表参数。
type SeasonListParams struct {
	Page     int
	PageSize int
	Order    string
	Sort     string
	Draft    int
}

// SeasonListData 合集列表响应。
type SeasonListData struct {
	Seasons  []SeasonEntry `json:"seasons"`
	Tip      SeasonTip     `json:"tip"`
	Total    int           `json:"total"`
	PlayType int           `json:"play_type"`
}

// SeasonEntry 合集列表项。
type SeasonEntry struct {
	Season       SeasonInfo      `json:"season"`
	Checkin      SeasonCheckin   `json:"checkin"`
	SeasonStat   SeasonStat      `json:"seasonStat"`
	Sections     SeasonSections  `json:"sections"`
	PartEpisodes []SeasonEpisode `json:"part_episodes"`
}

// SeasonInfo 合集信息。
type SeasonInfo struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Desc           string `json:"desc"`
	Cover          string `json:"cover"`
	IsEnd          int    `json:"isEnd"`
	Mid            int64  `json:"mid"`
	IsAct          int    `json:"isAct"`
	IsPay          int    `json:"is_pay"`
	State          int    `json:"state"`
	PartState      int    `json:"partState"`
	SignState      int    `json:"signState"`
	RejectReason   string `json:"rejectReason"`
	CTime          int64  `json:"ctime"`
	MTime          int64  `json:"mtime"`
	NoSection      int    `json:"no_section"`
	Forbid         int    `json:"forbid"`
	ProtocolID     string `json:"protocol_id"`
	EPNum          int    `json:"ep_num"`
	SeasonPrice    int    `json:"season_price"`
	IsOpened       int    `json:"is_opened"`
	HasChargingPay int    `json:"has_charging_pay"`
}

// SeasonTip 合集列表提示信息。
type SeasonTip struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// SeasonCheckin 合集审核信息。
type SeasonCheckin struct {
	Status       int    `json:"status"`
	StatusReason string `json:"status_reason"`
	SeasonStatus int    `json:"season_status"`
}

// SeasonStat 合集统计信息。
type SeasonStat struct {
	View         int `json:"view"`
	Danmaku      int `json:"danmaku"`
	Reply        int `json:"reply"`
	Fav          int `json:"fav"`
	Coin         int `json:"coin"`
	Share        int `json:"share"`
	NowRank      int `json:"nowRank"`
	HisRank      int `json:"hisRank"`
	Like         int `json:"like"`
	Subscription int `json:"subscription"`
	VT           int `json:"vt"`
}

// SeasonSections 合集小节列表。
type SeasonSections struct {
	Sections []SeasonSection `json:"sections"`
}

// SeasonSection 合集小节信息。
type SeasonSection struct {
	ID             int64       `json:"id"`
	Type           int         `json:"type"`
	SeasonID       int64       `json:"seasonId"`
	Title          string      `json:"title"`
	Order          int         `json:"order"`
	IsEnd          int         `json:"isEnd,omitempty"`
	State          int         `json:"state"`
	PartState      int         `json:"partState"`
	RejectReason   string      `json:"rejectReason"`
	CTime          int64       `json:"ctime"`
	MTime          int64       `json:"mtime"`
	EPCount        int         `json:"epCount"`
	Cover          string      `json:"cover"`
	Show           int         `json:"show,omitempty"`
	HasChargingPay int         `json:"has_charging_pay"`
	HasPUGVPay     int         `json:"has_pugv_pay,omitempty"`
	Episodes       interface{} `json:"Episodes,omitempty"`
}

// SeasonEpisode 合集中的视频信息。
type SeasonEpisode struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	AID          int64  `json:"aid"`
	BVID         string `json:"bvid"`
	CID          int64  `json:"cid"`
	SeasonID     int64  `json:"seasonId"`
	SectionID    int64  `json:"sectionId"`
	Order        int    `json:"order"`
	VideoTitle   string `json:"videoTitle"`
	ArchiveTitle string `json:"archiveTitle"`
	ArchiveState int    `json:"archiveState"`
	RejectReason string `json:"rejectReason"`
	State        int    `json:"state"`
	Cover        string `json:"cover"`
	IsFree       int    `json:"is_free"`
	AIDOwner     bool   `json:"aid_owner"`
	ChargingPay  int    `json:"charging_pay"`
	MemberFirst  int    `json:"member_first,omitempty"`
	PUGVPay      int    `json:"pugv_pay,omitempty"`
}

// SeasonCreateRequest 创建合集请求。
type SeasonCreateRequest struct {
	Title       string
	Desc        string
	Cover       string
	SeasonPrice int
}

// SeasonEpisodeInput 添加合集视频请求项。
type SeasonEpisodeInput struct {
	Title       string `json:"title"`
	AID         int64  `json:"aid"`
	CID         int64  `json:"cid"`
	ChargingPay int    `json:"charging_pay,omitempty"`
}

// SeasonAddEpisodesRequest 添加视频到合集请求。
type SeasonAddEpisodesRequest struct {
	SectionID int64
	Episodes  []SeasonEpisodeInput
}

// SeasonSectionEditRequest 编辑合集小节请求。
type SeasonSectionEditRequest struct {
	Section SeasonSectionEditInfo
	Sorts   []SeasonEpisodeOrder
}

// SeasonSectionEditInfo 编辑合集小节信息。
type SeasonSectionEditInfo struct {
	ID       int64  `json:"id"`
	Type     int    `json:"type"`
	SeasonID int64  `json:"seasonId"`
	Title    string `json:"title"`
}

// SeasonEpisodeOrder 小节内视频排序。
type SeasonEpisodeOrder struct {
	ID    int64 `json:"id"`
	Order int   `json:"order"`
}

// SeasonEditRequest 编辑合集请求。
type SeasonEditRequest struct {
	Season SeasonEditInfo
	Sorts  []SeasonSectionOrder
}

// SeasonEditInfo 编辑合集信息。
type SeasonEditInfo struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	Desc        string `json:"desc,omitempty"`
	SeasonPrice int    `json:"season_price,omitempty"`
	IsEnd       int    `json:"isEnd,omitempty"`
}

// SeasonSectionOrder 合集小节排序。
type SeasonSectionOrder struct {
	ID   int64 `json:"id"`
	Sort int   `json:"sort"`
}

// SeasonSectionDetail 合集小节详情。
type SeasonSectionDetail struct {
	Section  SeasonSection   `json:"section"`
	Episodes []SeasonEpisode `json:"episodes"`
}

// GetSeasonList 获取合集列表。
func (c *Client) GetSeasonList(params *SeasonListParams, cookies string) (*SeasonListData, error) {
	query := url.Values{}
	query.Set("pn", fmt.Sprintf("%d", seasonListPage(params)))
	query.Set("ps", fmt.Sprintf("%d", seasonListPageSize(params)))

	if params != nil {
		if value := strings.TrimSpace(params.Order); value != "" {
			query.Set("order", value)
		}
		if value := strings.TrimSpace(params.Sort); value != "" {
			query.Set("sort", value)
		}
		if params.Draft > 0 {
			query.Set("draft", fmt.Sprintf("%d", params.Draft))
		}
	}

	req, err := c.newCreativeCenterRequest(http.MethodGet, "https://member.bilibili.com/x2/creative/web/seasons?"+query.Encode(), nil, cookies)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	data, err := doCreativeCenterRequest[SeasonListData](c, req, "get season list")
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CreateSeason 创建合集。
func (c *Client) CreateSeason(createReq *SeasonCreateRequest, cookies string) (int64, error) {
	if createReq == nil {
		return 0, fmt.Errorf("create request is required")
	}
	if strings.TrimSpace(createReq.Title) == "" {
		return 0, fmt.Errorf("title is required")
	}
	if strings.TrimSpace(createReq.Cover) == "" {
		return 0, fmt.Errorf("cover is required")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return 0, err
	}

	formData := url.Values{}
	formData.Set("title", createReq.Title)
	formData.Set("desc", createReq.Desc)
	formData.Set("cover", createReq.Cover)
	formData.Set("season_price", fmt.Sprintf("%d", createReq.SeasonPrice))
	formData.Set("csrf", csrf)

	req, err := c.newCreativeCenterRequest(http.MethodPost, "https://member.bilibili.com/x2/creative/web/season/add", strings.NewReader(formData.Encode()), cookies)
	if err != nil {
		return 0, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	seasonID, err := doCreativeCenterRequest[int64](c, req, "create season")
	if err != nil {
		return 0, err
	}

	return *seasonID, nil
}

// AddSeasonEpisodes 添加视频到合集小节。
func (c *Client) AddSeasonEpisodes(addReq *SeasonAddEpisodesRequest, cookies string) error {
	if addReq == nil {
		return fmt.Errorf("add request is required")
	}
	if addReq.SectionID <= 0 {
		return fmt.Errorf("section ID is required")
	}
	if len(addReq.Episodes) == 0 {
		return fmt.Errorf("episodes are required")
	}
	for index, episode := range addReq.Episodes {
		if episode.AID <= 0 {
			return fmt.Errorf("episode[%d] aid is required", index)
		}
		if episode.CID <= 0 {
			return fmt.Errorf("episode[%d] cid is required", index)
		}
		if strings.TrimSpace(episode.Title) == "" {
			return fmt.Errorf("episode[%d] title is required", index)
		}
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"sectionId":  addReq.SectionID,
		"section_id": addReq.SectionID,
		"episodes":   addReq.Episodes,
		"episode":    addReq.Episodes,
	}

	if err := c.postCreativeCenterJSONWithCSRF("https://member.bilibili.com/x2/creative/web/season/section/episodes/add", payload, cookies, csrf, "add season episodes"); err != nil {
		return err
	}

	return nil
}

// EditSeasonSection 编辑合集小节。
func (c *Client) EditSeasonSection(editReq *SeasonSectionEditRequest, cookies string) error {
	if editReq == nil {
		return fmt.Errorf("edit request is required")
	}
	if editReq.Section.ID <= 0 {
		return fmt.Errorf("section ID is required")
	}
	if editReq.Section.SeasonID <= 0 {
		return fmt.Errorf("season ID is required")
	}
	if strings.TrimSpace(editReq.Section.Title) == "" {
		return fmt.Errorf("section title is required")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return err
	}

	sectionType := editReq.Section.Type
	if sectionType == 0 {
		sectionType = 1
	}

	sorts := make([]map[string]interface{}, 0, len(editReq.Sorts))
	for index, sortItem := range editReq.Sorts {
		if sortItem.ID <= 0 {
			return fmt.Errorf("sorts[%d] id is required", index)
		}
		sorts = append(sorts, map[string]interface{}{
			"id":    sortItem.ID,
			"order": sortItem.Order,
			"sort":  sortItem.Order,
		})
	}

	payload := map[string]interface{}{
		"section": map[string]interface{}{
			"id":       editReq.Section.ID,
			"type":     sectionType,
			"seasonId": editReq.Section.SeasonID,
			"title":    editReq.Section.Title,
		},
		"sorts": sorts,
	}

	if err := c.postCreativeCenterJSONWithCSRF("https://member.bilibili.com/x2/creative/web/season/section/edit", payload, cookies, csrf, "edit season section"); err != nil {
		return err
	}

	return nil
}

// EditSeason 编辑合集信息。
func (c *Client) EditSeason(editReq *SeasonEditRequest, cookies string) error {
	if editReq == nil {
		return fmt.Errorf("edit request is required")
	}
	if editReq.Season.ID <= 0 {
		return fmt.Errorf("season ID is required")
	}
	if strings.TrimSpace(editReq.Season.Title) == "" {
		return fmt.Errorf("season title is required")
	}
	if strings.TrimSpace(editReq.Season.Cover) == "" {
		return fmt.Errorf("season cover is required")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return err
	}

	sorts := make([]SeasonSectionOrder, 0, len(editReq.Sorts))
	for index, sortItem := range editReq.Sorts {
		if sortItem.ID <= 0 {
			return fmt.Errorf("sorts[%d] id is required", index)
		}
		sorts = append(sorts, sortItem)
	}

	payload := map[string]interface{}{
		"season": editReq.Season,
		"sorts":  sorts,
	}

	if err := c.postCreativeCenterJSONWithCSRF("https://member.bilibili.com/x2/creative/web/season/edit", payload, cookies, csrf, "edit season"); err != nil {
		return err
	}

	return nil
}

// DeleteSeason 删除合集。
func (c *Client) DeleteSeason(seasonID int64, cookies string) error {
	if seasonID <= 0 {
		return fmt.Errorf("season ID is required")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return err
	}

	formData := url.Values{}
	formData.Set("id", fmt.Sprintf("%d", seasonID))
	formData.Set("csrf", csrf)

	req, err := c.newCreativeCenterRequest(http.MethodPost, "https://member.bilibili.com/x2/creative/web/season/del", strings.NewReader(formData.Encode()), cookies)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = doCreativeCenterRequest[struct{}](c, req, "delete season")
	return err
}

// GetSeasonSection 获取合集小节中的视频。
func (c *Client) GetSeasonSection(seasonID int64, cookies string) (*SeasonSectionDetail, error) {
	if seasonID <= 0 {
		return nil, fmt.Errorf("season ID is required")
	}

	apiURL := fmt.Sprintf("https://member.bilibili.com/x2/creative/web/season/section?id=%d", seasonID)
	req, err := c.newCreativeCenterRequest(http.MethodGet, apiURL, nil, cookies)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	data, err := doCreativeCenterRequest[SeasonSectionDetail](c, req, "get season section")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) postCreativeCenterJSONWithCSRF(apiURL string, payload interface{}, cookies string, csrf string, action string) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	queryURL := apiURL + "?csrf=" + url.QueryEscape(csrf)
	req, err := c.newCreativeCenterRequest(http.MethodPost, queryURL, bytes.NewReader(body), cookies)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = doCreativeCenterRequest[struct{}](c, req, action)
	return err
}

func (c *Client) newCreativeCenterRequest(method string, apiURL string, body io.Reader, cookies string) (*http.Request, error) {
	req, err := http.NewRequest(method, apiURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
	if strings.TrimSpace(cookies) != "" {
		req.Header.Set("Cookie", cookies)
	}

	return req, nil
}

func doCreativeCenterRequest[T any](c *Client, req *http.Request, action string) (*T, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result creativeCenterResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("%s failed: code=%d, message=%s", action, result.Code, result.Message)
	}

	return &result.Data, nil
}

func requireCSRFFromCookies(cookies string) (string, error) {
	csrf := strings.TrimSpace(extractCSRFFromCookies(cookies))
	if csrf == "" {
		return "", fmt.Errorf("csrf token not found in cookies")
	}
	return csrf, nil
}

func seasonListPage(params *SeasonListParams) int {
	if params == nil || params.Page <= 0 {
		return 1
	}
	return params.Page
}

func seasonListPageSize(params *SeasonListParams) int {
	if params == nil || params.PageSize <= 0 {
		return 30
	}
	return params.PageSize
}
