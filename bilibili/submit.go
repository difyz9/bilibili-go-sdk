package bilibili

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// SubmitVideo 提交视频到B站，使用Web投稿接口提交新稿件
func (uc *UploadClient) SubmitVideo(studio *Studio) (*ResponseData, error) {
	return uc.SubmitVideoWeb(studio)
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
	return uc.uploadCoverDataURI(dataURI)
}

func (uc *UploadClient) uploadCoverDataURI(dataURI string) (string, error) {
	if strings.TrimSpace(dataURI) == "" {
		return "", fmt.Errorf("cover data is empty")
	}

	// 获取CSRF token
	csrf, err := uc.loginInfo.GetCSRFToken()
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("ts", fmt.Sprintf("%d", time.Now().UnixMilli()))

	form := url.Values{}
	form.Set("csrf", csrf)
	form.Set("cover", dataURI)

	req, err := uc.newCreativeCenterRequest(
		http.MethodPost,
		"/x/vu/web/cover/up",
		query,
		bytes.NewBufferString(form.Encode()),
		"application/x-www-form-urlencoded",
	)
	if err != nil {
		return "", err
	}

	var result ResponseData
	if err := uc.doJSONRequest(req, &result); err != nil {
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
	return uc.uploadCoverDataURI(dataURI)
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