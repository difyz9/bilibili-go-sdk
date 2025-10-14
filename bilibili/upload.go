package bilibili

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadClient 上传客户端
type UploadClient struct {
	client       *Client
	uploadClient *http.Client // 专门用于上传的 HTTP 客户端（更长的超时时间）
	loginInfo    *LoginInfo
}

// NewUploadClient 创建上传客户端
func NewUploadClient(loginInfo *LoginInfo, opts ...Option) *UploadClient {
	config := DefaultConfig()
	config.ApplyOptions(opts...)

	// 创建专门用于上传的客户端，优化网络参数
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  true, // 禁用压缩以减少CPU使用
	}
	
	uploadClient := &http.Client{
		Timeout:   15 * time.Minute, // 增加超时时间到15分钟
		Transport: transport,
	}

	return &UploadClient{
		client:       NewClient(opts...),
		uploadClient: uploadClient,
		loginInfo:    loginInfo,
	}
}

// retryFunc 重试函数，针对网络错误进行优化
func retryFunc(fn func() error) error {
	maxRetries := 5 // 增加重试次数
	wait := 2.0     // 增加初始等待时间

	for retries := maxRetries; retries > 0; retries-- {
		err := fn()
		if err == nil {
			return nil
		}

		// 检查是否是可重试的网络错误
		if retries > 1 && IsNetworkError(err) {
			// 指数退避算法 + 随机抖动
			jitter := rand.Float64() * 2.0 // 增加抖动范围
			waitTime := math.Min(jitter+wait, 120.0) // 增加最大等待时间
			
			log.Printf("🔄 Retry attempt #%d/%d. Network error detected, sleeping %.2fs before retry. Error: %v", 
				maxRetries-retries+1, maxRetries, waitTime, err)
			
			time.Sleep(time.Duration(waitTime * float64(time.Second)))
			wait *= 1.8 // 稍微降低指数增长率
		} else if retries > 1 {
			// 非网络错误，快速重试
			log.Printf("⚡ Quick retry #%d/%d for non-network error: %v", maxRetries-retries+1, maxRetries, err)
			time.Sleep(1 * time.Second)
		} else {
			return err
		}
	}
	return nil
}

// UploadVideo 上传视频文件
func (uc *UploadClient) UploadVideo(videoPath string) (*Video, error) {
	// 1. 获取文件信息
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	fileName := filepath.Base(videoPath)
	fileSize := fileInfo.Size()

	log.Printf("Uploading video: %s, size: %d bytes", fileName, fileSize)

	// 2. 预上传获取配置
	preUploadInfo, err := uc.preUpload(fileName, fileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-upload: %v", err)
	}

	// 3. 获取上传ID
	uploadID, err := uc.getUploadID(preUploadInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload ID: %v", err)
	}

	log.Printf("Got upload ID: %s", uploadID)

	// 4. 分块上传文件
	parts, err := uc.uploadChunks(videoPath, preUploadInfo, uploadID)
	if err != nil {
		return nil, fmt.Errorf("failed to upload chunks: %v", err)
	}

	log.Printf("Uploaded %d chunks", len(parts))

	// 5. 完成上传
	video, err := uc.completeUpload(preUploadInfo, uploadID, parts, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to complete upload: %v", err)
	}

	// 设置 title 为原始文件名（去除扩展名）
	titleWithoutExt := fileName
	if ext := filepath.Ext(fileName); ext != "" {
		titleWithoutExt = strings.TrimSuffix(fileName, ext)
	}
	video.Title = titleWithoutExt

	log.Printf("Upload completed. Video filename: %s, title: %s", video.Filename, video.Title)

	return video, nil
}

// UploadVideoFromURL 从 URL 上传视频文件到 Bilibili（不使用临时文件）
func (uc *UploadClient) UploadVideoFromURL(videoURL, fileName string, fileSize int64) (*Video, error) {
	log.Printf("Uploading video from URL: %s, filename: %s, size: %d bytes", videoURL, fileName, fileSize)

	// 1. 预上传获取配置
	preUploadInfo, err := uc.preUpload(fileName, fileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-upload: %v", err)
	}

	// 2. 获取上传ID
	uploadID, err := uc.getUploadID(preUploadInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload ID: %v", err)
	}

	log.Printf("Got upload ID: %s", uploadID)

	// 3. 分块上传文件（从URL流式上传）
	parts, err := uc.uploadChunksFromURL(videoURL, preUploadInfo, uploadID, fileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to upload chunks from URL: %v", err)
	}

	log.Printf("Uploaded %d chunks from URL", len(parts))

	// 4. 完成上传
	video, err := uc.completeUpload(preUploadInfo, uploadID, parts, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to complete upload: %v", err)
	}

	// 设置 title 为原始文件名（去除扩展名）
	titleWithoutExt := fileName
	if ext := filepath.Ext(fileName); ext != "" {
		titleWithoutExt = strings.TrimSuffix(fileName, ext)
	}
	video.Title = titleWithoutExt

	log.Printf("Upload from URL completed. Video filename: %s, title: %s", video.Filename, video.Title)

	return video, nil
}