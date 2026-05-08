package bilibili

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type apiResponse[T any] struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data"`
	TTL       int    `json:"ttl,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// UploadTemplate 创作中心投稿模板
type UploadTemplate struct {
	TID         int64  `json:"tid"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Tags        string `json:"tags"`
	Description string `json:"description"`
	TypeID      int    `json:"typeid"`
	Copyright   int    `json:"copyright"`
	Attribute   int    `json:"attribute"`
	IsDefault   int    `json:"is_default"`
}

// UpdateUploadTemplateRequest 编辑模板请求
type UpdateUploadTemplateRequest struct {
	TID         int64  `json:"tid"`
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Description string `json:"description,omitempty"`
	TypeID      int    `json:"typeid,omitempty"`
	ArcType     string `json:"arctype,omitempty"`
	IsDefault   *int   `json:"is_default,omitempty"`
}

// TopicQueryRequest 查询话题请求
type TopicQueryRequest struct {
	TypeID int
	PN     int
	PS     int
	Title  string
}

// TopicInfo 话题信息
type TopicInfo struct {
	TopicID             int64  `json:"topic_id"`
	TopicName           string `json:"topic_name"`
	Description         string `json:"description"`
	MissionID           int64  `json:"mission_id"`
	ActivityText        string `json:"activity_text"`
	ActivityDescription string `json:"activity_description"`
}

// TopicSearchRequest 话题搜索请求
type TopicSearchRequest struct {
	PageSize int
	Offset   int
	Keywords string
}

// TopicSearchPageInfo 话题搜索分页信息
type TopicSearchPageInfo struct {
	HasMore    bool `json:"has_more"`
	Offset     int  `json:"offset"`
	PageNumber int  `json:"page_number"`
}

// TopicSearchItem 话题搜索结果项
type TopicSearchItem struct {
	ActProtocol  string `json:"act_protocol"`
	ActivitySign string `json:"activity_sign"`
	Description  string `json:"description"`
	ID           int64  `json:"id"`
	MissionID    int64  `json:"mission_id"`
	Name         string `json:"name"`
	State        int    `json:"state"`
	Uname        string `json:"uname"`
}

// TopicSearchResult 话题搜索结果
type TopicSearchResult struct {
	HasCreateJurisdiction bool                `json:"has_create_jurisdiction"`
	IsNewTopic            bool                `json:"is_new_topic"`
	Tips                  string              `json:"tips"`
	PageInfo              TopicSearchPageInfo `json:"page_info"`
	Topics                []TopicSearchItem   `json:"topics"`
}

// ArchiveDescFormat 简介格式信息
type ArchiveDescFormat struct {
	TypeID    int    `json:"typeid"`
	ID        int64  `json:"id"`
	Lang      int    `json:"lang"`
	Copyright int    `json:"copyright"`
	Components string `json:"components"`
}

// UploadLine 上传线路
type UploadLine struct {
	OS       string `json:"os"`
	Query    string `json:"query"`
	ProbeURL string `json:"probe_url"`
}

// UploadProbe 上传探测信息
type UploadProbe struct {
	Post float64 `json:"post"`
}

// UploadLineProbeResponse 获取上传线路响应
type UploadLineProbeResponse struct {
	OK    int          `json:"OK"`
	Lines []UploadLine `json:"lines"`
	Probe UploadProbe  `json:"probe"`
}

// ArchiveTypePredictRequest 稿件类型预测请求
type ArchiveTypePredictRequest struct {
	Filename string
	Title    string
	UploadID string
}

// HumanTypeRef 新分区引用
type HumanTypeRef struct {
	ID int `json:"id"`
}

// ArchiveTypePrediction 稿件类型预测结果
type ArchiveTypePrediction struct {
	ID            int           `json:"id"`
	Parent        int           `json:"parent"`
	ParentName    string        `json:"parent_name"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Desc          string        `json:"desc"`
	IntroOriginal string        `json:"intro_original"`
	IntroCopy     string        `json:"intro_copy"`
	Notice        string        `json:"notice"`
	CopyRight     int           `json:"copy_right"`
	Show          bool          `json:"show"`
	Rank          int           `json:"rank"`
	MaxVideoCount int           `json:"max_video_count"`
	RequestID     string        `json:"request_id"`
	HumanType     *HumanTypeRef `json:"human_type"`
}

// HumanTypeInfo 新分区信息
type HumanTypeInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WebVideoEditRequest Web 编辑稿件请求
type WebVideoEditRequest struct {
	AID        int64 `json:"aid"`
	Studio
	NewWebEdit int `json:"new_web_edit,omitempty"`
}

func (uc *UploadClient) newCreativeCenterRequest(method, path string, query url.Values, body io.Reader, contentType string) (*http.Request, error) {
	requestURL := "https://member.bilibili.com" + path
	if len(query) > 0 {
		requestURL += "?" + query.Encode()
	}

	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", uc.client.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
	if uc.loginInfo != nil {
		cookies := uc.loginInfo.GetCookieString()
		if cookies != "" {
			req.Header.Set("Cookie", cookies)
		}
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

func (uc *UploadClient) doJSONRequest(req *http.Request, out any) error {
	resp, err := uc.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("request failed: status=%s body=%s", resp.Status, string(body))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode response failed: %w, body: %s", err, string(body))
	}

	return nil
}

func normalizeStudioForSubmitWeb(studio Studio) Studio {
	if studio.DescFormatId == 0 {
		studio.DescFormatId = 9999
	}
	if studio.Recreate == 0 {
		studio.Recreate = -1
	}
	if studio.WebOS == 0 {
		studio.WebOS = 3
	}
	if studio.Is360 == 0 {
		studio.Is360 = -1
	}
	if studio.Subtitle == (Subtitle{}) {
		studio.Subtitle = Subtitle{Open: 0, Lan: ""}
	}
	for index := range studio.Videos {
		if studio.Videos[index].Title == "" {
			studio.Videos[index].Title = studio.Title
		}
	}
	return studio
}

func normalizeStudioForEditWeb(studio Studio) Studio {
	studio = normalizeStudioForSubmitWeb(studio)
	if studio.WebOS == 3 {
		studio.WebOS = 1
	}
	return studio
}

// SubmitVideoWeb 使用Web接口提交视频稿件
func (uc *UploadClient) SubmitVideoWeb(studio *Studio) (*ResponseData, error) {
	if studio == nil {
		return nil, fmt.Errorf("studio is nil")
	}
	if len(studio.Videos) == 0 {
		return nil, fmt.Errorf("no videos provided")
	}

	normalized := normalizeStudioForSubmitWeb(*studio)
	for index, video := range normalized.Videos {
		if video.Filename == "" {
			return nil, fmt.Errorf("video[%d] filename is required", index)
		}
		if video.CID <= 0 {
			return nil, fmt.Errorf("video[%d] cid is required for web submit", index)
		}
	}

	csrf, err := uc.loginInfo.GetCSRFToken()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("csrf", csrf)
	query.Set("ts", strconv.FormatInt(time.Now().UnixMilli(), 10))

	req, err := uc.newCreativeCenterRequest(http.MethodPost, "/x/vu/web/add/v3", query, bytes.NewReader(body), "application/json; charset=utf-8")
	if err != nil {
		return nil, err
	}

	var result ResponseData
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("submit video failed: code=%d, message=%s", result.Code, result.Message)
	}

	return &result, nil
}

// EditVideoWeb 使用Web接口编辑视频稿件
func (uc *UploadClient) EditVideoWeb(request *WebVideoEditRequest) (*ResponseData, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if request.AID <= 0 {
		return nil, fmt.Errorf("aid is required")
	}
	if len(request.Videos) == 0 {
		return nil, fmt.Errorf("no videos provided")
	}

	payload := *request
	payload.Studio = normalizeStudioForEditWeb(request.Studio)
	for index, video := range payload.Videos {
		if video.CID <= 0 {
			return nil, fmt.Errorf("video[%d] cid is required for web edit", index)
		}
	}

	csrf, err := uc.loginInfo.GetCSRFToken()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("csrf", csrf)
	query.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))

	req, err := uc.newCreativeCenterRequest(http.MethodPost, "/x/vu/web/edit", query, bytes.NewReader(body), "application/json; charset=utf-8")
	if err != nil {
		return nil, err
	}

	var result ResponseData
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("edit video failed: code=%d, message=%s", result.Code, result.Message)
	}

	return &result, nil
}

// GetUploadTemplates 获取上传模板列表
func (uc *UploadClient) GetUploadTemplates() ([]UploadTemplate, error) {
	query := url.Values{}
	query.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))

	req, err := uc.newCreativeCenterRequest(http.MethodGet, "/x/vupre/web/tpls", query, nil, "")
	if err != nil {
		return nil, err
	}

	var result apiResponse[[]UploadTemplate]
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("get upload templates failed: code=%d, message=%s", result.Code, result.Message)
	}

	return result.Data, nil
}

// UpdateUploadTemplate 编辑上传模板
func (uc *UploadClient) UpdateUploadTemplate(request *UpdateUploadTemplateRequest) error {
	if request == nil {
		return fmt.Errorf("request is nil")
	}
	if request.TID <= 0 {
		return fmt.Errorf("tid is required")
	}

	csrf, err := uc.loginInfo.GetCSRFToken()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"tid":  request.TID,
		"csrf": csrf,
	}
	if request.Name != "" {
		payload["name"] = request.Name
	}
	if request.Title != "" {
		payload["title"] = request.Title
	}
	if request.Keywords != "" {
		payload["keywords"] = request.Keywords
	}
	if request.Description != "" {
		payload["description"] = request.Description
	}
	if request.TypeID > 0 {
		payload["typeid"] = request.TypeID
	}
	if request.ArcType != "" {
		payload["arctype"] = request.ArcType
	}
	if request.IsDefault != nil {
		payload["is_default"] = *request.IsDefault
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	query := url.Values{}
	query.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))

	req, err := uc.newCreativeCenterRequest(http.MethodPost, "/x/vupre/web/tpl/update", query, bytes.NewReader(body), "application/json; charset=utf-8")
	if err != nil {
		return err
	}

	var result ResponseData
	if err := uc.doJSONRequest(req, &result); err != nil {
		return err
	}
	if result.Code != 0 {
		return fmt.Errorf("update upload template failed: code=%d, message=%s", result.Code, result.Message)
	}

	return nil
}

// QueryTopics 查询话题
func (uc *UploadClient) QueryTopics(request *TopicQueryRequest) ([]TopicInfo, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if request.PS <= 0 {
		return nil, fmt.Errorf("ps must be greater than 0")
	}

	query := url.Values{}
	query.Set("pn", strconv.Itoa(request.PN))
	query.Set("ps", strconv.Itoa(request.PS))
	query.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))
	if request.TypeID > 0 {
		query.Set("type_id", strconv.Itoa(request.TypeID))
	}
	if request.Title != "" {
		query.Set("title", request.Title)
	}

	req, err := uc.newCreativeCenterRequest(http.MethodGet, "/x/vupre/web/topic/type", query, nil, "")
	if err != nil {
		return nil, err
	}

	var result apiResponse[[]TopicInfo]
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("query topics failed: code=%d, message=%s", result.Code, result.Message)
	}

	return result.Data, nil
}

// SearchTopics 搜索话题
func (uc *UploadClient) SearchTopics(request *TopicSearchRequest) (*TopicSearchResult, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	query := url.Values{}
	query.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))
	if request.PageSize > 0 {
		query.Set("page_size", strconv.Itoa(request.PageSize))
	}
	if request.Offset > 0 {
		query.Set("offset", strconv.Itoa(request.Offset))
	}
	if request.Keywords != "" {
		query.Set("keywords", request.Keywords)
	}

	req, err := uc.newCreativeCenterRequest(http.MethodGet, "/x/vupre/web/topic/search", query, nil, "")
	if err != nil {
		return nil, err
	}

	var result apiResponse[struct {
		Result TopicSearchResult `json:"result"`
	}]
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("search topics failed: code=%d, message=%s", result.Code, result.Message)
	}

	return &result.Data.Result, nil
}

// GetArchiveDescFormat 获取简介相关信息
func (uc *UploadClient) GetArchiveDescFormat(typeID, copyright int) (*ArchiveDescFormat, error) {
	if typeID <= 0 {
		return nil, fmt.Errorf("typeid is required")
	}
	if copyright <= 0 {
		return nil, fmt.Errorf("copyright is required")
	}

	query := url.Values{}
	query.Set("typeid", strconv.Itoa(typeID))
	query.Set("copyright", strconv.Itoa(copyright))
	query.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))

	req, err := uc.newCreativeCenterRequest(http.MethodGet, "/x/vupre/web/archive/desc/format", query, nil, "")
	if err != nil {
		return nil, err
	}

	var result apiResponse[*ArchiveDescFormat]
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("get archive desc format failed: code=%d, message=%s", result.Code, result.Message)
	}

	return result.Data, nil
}

// ProbeUploadLines 获取上传线路
func (uc *UploadClient) ProbeUploadLines() (*UploadLineProbeResponse, error) {
	req, err := http.NewRequest(http.MethodGet, "https://member.bilibili.com/preupload?r=probe", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", uc.client.userAgent)

	var result UploadLineProbeResponse
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.OK != 1 {
		return nil, fmt.Errorf("probe upload lines failed: ok=%d", result.OK)
	}

	return &result, nil
}

// PredictArchiveTypes 预测稿件类型
func (uc *UploadClient) PredictArchiveTypes(request *ArchiveTypePredictRequest) ([]ArchiveTypePrediction, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}

	csrf, err := uc.loginInfo.GetCSRFToken()
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("csrf", csrf)
	query.Set("ts", strconv.FormatInt(time.Now().UnixMilli(), 10))

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("filename", request.Filename); err != nil {
		return nil, err
	}
	if request.Title != "" {
		if err := writer.WriteField("title", request.Title); err != nil {
			return nil, err
		}
	}
	if request.UploadID != "" {
		if err := writer.WriteField("upload_id", request.UploadID); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := uc.newCreativeCenterRequest(http.MethodPost, "/x/vupre/web/archive/types/predict", query, &body, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var result apiResponse[[]ArchiveTypePrediction]
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("predict archive types failed: code=%d, message=%s", result.Code, result.Message)
	}

	for index := range result.Data {
		result.Data[index].RequestID = result.RequestID
	}

	return result.Data, nil
}

// GetHumanTypeList 获取新分区列表
func (uc *UploadClient) GetHumanTypeList() ([]HumanTypeInfo, error) {
	query := url.Values{}
	query.Set("t", strconv.FormatInt(time.Now().UnixMilli(), 10))

	req, err := uc.newCreativeCenterRequest(http.MethodGet, "/x/vupre/web/archive/human/type2/list", query, nil, "")
	if err != nil {
		return nil, err
	}

	var result struct {
		Code     int             `json:"code"`
		Message  string          `json:"message"`
		TTL      int             `json:"ttl"`
		TypeList []HumanTypeInfo `json:"type_list"`
	}
	if err := uc.doJSONRequest(req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("get human type list failed: code=%d, message=%s", result.Code, result.Message)
	}

	return result.TypeList, nil
}