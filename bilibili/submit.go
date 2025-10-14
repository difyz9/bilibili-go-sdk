package bilibili

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// SubmitVideo 提交视频到B站，使用APP接口提交新稿件
func (uc *UploadClient) SubmitVideo(studio *Studio) (*ResponseData, error) {
	log.Println("=== 提交视频到B站 ===")

	// 验证必填字段
	if len(studio.Videos) == 0 {
		return nil, fmt.Errorf("no videos provided")
	}

	// 确保 videos 中的 filename 格式正确（不应包含扩展名）
	for i, video := range studio.Videos {
		log.Printf("Video[%d]: filename=%s, title=%s", i, video.Filename, video.Title)
		// 如果 filename 包含路径分隔符或扩展名，进行清理
		if strings.Contains(video.Filename, "/") || strings.Contains(video.Filename, "\\") {
			log.Printf("Warning: Video filename contains path separator: %s", video.Filename)
		}
	}

	// 构建APP接口提交参数
	params := url.Values{}
	params.Set("access_key", uc.loginInfo.TokenInfo.AccessToken)
	params.Set("appkey", BiliTVAppKey)
	params.Set("build", "7800300")
	params.Set("c_locale", "zh-Hans_CN")
	params.Set("channel", "bili")
	params.Set("disable_rcmd", "0")
	params.Set("mobi_app", "android")
	params.Set("platform", "android")
	params.Set("s_locale", "zh-Hans_CN")
	params.Set("statistics", `"appId":1,"platform":3,"version":"7.80.0","abtest":""`)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))

	// 计算签名
	urlencoded := params.Encode()
	sign := Sign(urlencoded, BiliTVAppSec)
	params.Set("sign", sign)

	// 打印详细的提交数据用于调试
	studioData, _ := json.MarshalIndent(studio, "", "  ")
	log.Printf("Submit studio data:\n%s", string(studioData))
	log.Printf("Submit URL params: %s", params.Encode())

	// 准备请求体
	jsonData, err := json.Marshal(studio)
	if err != nil {
		return nil, err
	}

	submitURL := "https://member.bilibili.com/x/vu/app/add?" + params.Encode()
	req, err := http.NewRequest("POST", submitURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 BiliDroid/7.80.0 (bbcallen@gmail.com) os/android model/MI 6 mobi_app/android build/7800300 channel/bili innerVer/7800310 osVer/13 network/2")

	resp, err := uc.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体进行调试
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Submit response status: %d", resp.StatusCode)
	log.Printf("Submit response body: %s", string(body))

	var result ResponseData
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	return &result, nil
}

// UploadCover 上传封面
func (uc *UploadClient) UploadCover(imagePath string) (string, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	// 转换为base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	dataURI := fmt.Sprintf("data:image/jpeg;base64,%s", base64Data)

	// 获取CSRF token
	csrf, err := uc.loginInfo.GetCSRFToken()
	if err != nil {
		return "", err
	}

	// 构建表单数据
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("cover", dataURI)
	writer.WriteField("csrf", csrf)
	writer.Close()

	req, err := http.NewRequest("POST", "https://member.bilibili.com/x/vu/web/cover/up", &buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", uc.loginInfo.GetCookieString())
	req.Header.Set("User-Agent", uc.client.userAgent)

	resp, err := uc.client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("upload cover failed: %s", result.Message)
	}

	if data, ok := result.Data.(map[string]interface{}); ok {
		if url, ok := data["url"].(string); ok {
			return url, nil
		}
	}

	return "", fmt.Errorf("failed to get cover URL from response")
}

// UploadCoverFromBytes 从字节数据上传封面
func (uc *UploadClient) UploadCoverFromBytes(imageData []byte, contentType string) (string, error) {
	// 转换为base64
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	dataURI := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)

	// 获取CSRF token
	csrf, err := uc.loginInfo.GetCSRFToken()
	if err != nil {
		return "", err
	}

	// 构建表单数据
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("cover", dataURI)
	writer.WriteField("csrf", csrf)
	writer.Close()

	req, err := http.NewRequest("POST", "https://member.bilibili.com/x/vu/web/cover/up", &buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Cookie", uc.loginInfo.GetCookieString())
	req.Header.Set("User-Agent", uc.client.userAgent)

	resp, err := uc.client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result ResponseData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("upload cover failed: %s", result.Message)
	}

	if data, ok := result.Data.(map[string]interface{}); ok {
		if url, ok := data["url"].(string); ok {
			return url, nil
		}
	}

	return "", fmt.Errorf("failed to get cover URL from response")
}

// UploadCoverFromURL 从URL上传封面
func (uc *UploadClient) UploadCoverFromURL(imageURL string) (string, error) {
	// 下载图片
	resp, err := uc.client.httpClient.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %v", err)
	}

	// 获取内容类型
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // 默认类型
	}

	return uc.UploadCoverFromBytes(imageData, contentType)
}