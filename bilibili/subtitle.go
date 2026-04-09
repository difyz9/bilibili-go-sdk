package bilibili

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// SubtitleLangZhCN is the Bilibili language code for Simplified Chinese subtitles.
	SubtitleLangZhCN = "zh-CN"
	// SubtitleLangZhTW is the Bilibili language code for Traditional Chinese subtitles.
	SubtitleLangZhTW = "zh-TW"
	// SubtitleLangEN is the Bilibili language code for generic English subtitles.
	SubtitleLangEN = "en"
	// SubtitleLangENUS is the Bilibili language code for US English subtitles.
	SubtitleLangENUS = "en-US"
)

// SubtitleUploader Bilibili字幕上传器
type SubtitleUploader struct {
	client    *Client
	loginInfo *LoginInfo
}

// SubtitleVideoInfo 字幕相关的视频信息结构
type SubtitleVideoInfo struct {
	CID int64 `json:"cid"`
	AID int64 `json:"aid"`
}

// SubtitleFile 字幕文件结构
type SubtitleFile struct {
	URL        string `json:"url"`
	Language   string `json:"lan"`
	SubtitleID int    `json:"subtitle_id"`
}

// SubtitleUploadResponse 字幕上传响应
type SubtitleUploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Location string `json:"location"`
		Etag     string `json:"etag"`
	} `json:"data"`
}

// SubtitleSaveResponse 字幕保存响应
type SubtitleSaveResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

var subtitleLanguageAliases = map[string]string{
	"zh":       SubtitleLangZhCN,
	"zh-cn":    SubtitleLangZhCN,
	"zh-hans":  SubtitleLangZhCN,
	"cmn":      SubtitleLangZhCN,
	"cmn-hans": SubtitleLangZhCN,
	"zh-tw":    SubtitleLangZhTW,
	"zh-hant":  SubtitleLangZhTW,
	"cmn-hant": SubtitleLangZhTW,
	"en":       SubtitleLangEN,
	"en-us":    SubtitleLangENUS,
}

// VideoInfoResponse 视频信息响应
type VideoInfoResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Archive struct {
			StateDesc string `json:"state_desc"`
		} `json:"archive"`
		Videos []struct {
			CID int64 `json:"cid"`
			AID int64 `json:"aid"`
		} `json:"videos"`
	} `json:"data"`
}

// NewSubtitleUploader 创建字幕上传器
func NewSubtitleUploader(client *Client, loginInfo *LoginInfo) *SubtitleUploader {
	return &SubtitleUploader{
		client:    client,
		loginInfo: loginInfo,
	}
}

// NormalizeSubtitleLanguage converts common app/workflow language tags to Bilibili-accepted subtitle language codes.
func NormalizeSubtitleLanguage(language string) string {
	normalized := subtitleLanguageAliases[strings.ToLower(strings.TrimSpace(language))]
	if normalized != "" {
		return normalized
	}
	return language
}

// GetVideoInfo 获取视频信息（CID和AID）
func (s *SubtitleUploader) GetVideoInfo(bvid string) (*SubtitleVideoInfo, error) {
	url := fmt.Sprintf("https://member.bilibili.com/x/vupre/web/archive/view?bvid=%s", bvid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 添加Cookie
	cookieStr := s.loginInfo.GetCookieString()
	req.Header.Set("Cookie", cookieStr)
	req.Header.Set("User-Agent", s.client.userAgent)

	// 重试机制
	var resp *http.Response
	for attempt := 0; attempt < 3; attempt++ {
		resp, err = s.client.httpClient.Do(req)
		if err == nil {
			break
		}
		if attempt < 2 {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	var response VideoInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("get video info failed: code=%d, message=%s", response.Code, response.Message)
	}

	if len(response.Data.Videos) == 0 {
		return nil, fmt.Errorf("video info is empty")
	}

	return &SubtitleVideoInfo{
		CID: response.Data.Videos[0].CID,
		AID: response.Data.Videos[0].AID,
	}, nil
}

// UploadSubtitleFile 上传字幕文件到Bilibili存储
func (s *SubtitleUploader) UploadSubtitleFile(subtitlePath string) (string, string, error) {
	// 获取CSRF Token
	csrf, err := s.loginInfo.GetCSRFToken()
	if err != nil {
		return "", "", fmt.Errorf("get CSRF token failed: %w", err)
	}

	// 创建multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加字段
	writer.WriteField("bucket", "subtitle")
	writer.WriteField("csrf", csrf)
	writer.WriteField("content_type", "application/x-subrip")

	// 添加文件
	file, err := os.Open(subtitlePath)
	if err != nil {
		return "", "", fmt.Errorf("open subtitle file failed: %w", err)
	}
	defer file.Close()

	fileWriter, err := writer.CreateFormFile("file", "subtitle.srt")
	if err != nil {
		return "", "", fmt.Errorf("create form file failed: %w", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", "", fmt.Errorf("copy file content failed: %w", err)
	}

	writer.Close()

	// 构建请求
	timestamp := time.Now().UnixMilli()
	url := fmt.Sprintf("https://api.bilibili.com/x/upload/web/image?t=%d&csrf=%s", timestamp, csrf)

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", "", fmt.Errorf("create upload request failed: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", s.loginInfo.GetCookieString())
	req.Header.Set("User-Agent", s.client.userAgent)

	// 发送请求
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read upload response failed: %w", err)
	}

	var response SubtitleUploadResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", "", fmt.Errorf("unmarshal upload response failed: %w", err)
	}

	if response.Code != 0 {
		return "", "", fmt.Errorf("upload subtitle file failed: code=%d, message=%s", response.Code, response.Message)
	}

	return response.Data.Location, response.Data.Etag, nil
}

// SaveSubtitleInfo 保存字幕信息到视频
func (s *SubtitleUploader) SaveSubtitleInfo(aid, cid int64, location, language string) error {
	language = NormalizeSubtitleLanguage(language)

	// 获取CSRF Token
	csrf, err := s.loginInfo.GetCSRFToken()
	if err != nil {
		return fmt.Errorf("get CSRF token failed: %w", err)
	}

	// 构建字幕文件信息
	subtitleFiles := []SubtitleFile{
		{
			URL:        location,
			Language:   language,
			SubtitleID: 0,
		},
	}

	filesJSON, err := json.Marshal(subtitleFiles)
	if err != nil {
		return fmt.Errorf("marshal subtitle files failed: %w", err)
	}

	// 创建multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("oid", strconv.FormatInt(cid, 10))
	writer.WriteField("type", "1")
	writer.WriteField("files", string(filesJSON))
	writer.WriteField("aid", strconv.FormatInt(aid, 10))
	writer.WriteField("csrf", csrf)

	writer.Close()

	// 构建请求
	timestamp := time.Now().UnixMilli()
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/dm/subtitle/draft/preSave?t=%d&csrf=%s", timestamp, csrf)

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return fmt.Errorf("create save request failed: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", s.loginInfo.GetCookieString())
	req.Header.Set("User-Agent", s.client.userAgent)

	// 发送请求
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("save request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read save response failed: %w", err)
	}

	var response SubtitleSaveResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("unmarshal save response failed: %w", err)
	}

	if response.Code != 0 {
		return fmt.Errorf("save subtitle info failed: code=%d, message=%s", response.Code, response.Message)
	}

	return nil
}

// UploadSubtitle 完整的字幕上传流程
func (s *SubtitleUploader) UploadSubtitle(bvid, subtitlePath, language string) error {
	language = NormalizeSubtitleLanguage(language)

	// 1. 获取视频信息
	videoInfo, err := s.GetVideoInfo(bvid)
	if err != nil {
		return fmt.Errorf("get video info failed: %w", err)
	}

	// 2. 上传字幕文件
	location, _, err := s.UploadSubtitleFile(subtitlePath)
	if err != nil {
		return fmt.Errorf("upload subtitle file failed: %w", err)
	}

	// 3. 保存字幕信息
	err = s.SaveSubtitleInfo(videoInfo.AID, videoInfo.CID, location, language)
	if err != nil {
		return fmt.Errorf("save subtitle info failed: %w", err)
	}

	return nil
}

// UploadSubtitle 客户端级别的字幕上传方法（便捷方法）
func (c *Client) UploadSubtitle(loginInfo *LoginInfo, bvid, subtitlePath, language string) error {
	uploader := NewSubtitleUploader(c, loginInfo)
	return uploader.UploadSubtitle(bvid, subtitlePath, language)
}
