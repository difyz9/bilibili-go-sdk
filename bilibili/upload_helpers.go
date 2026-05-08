package bilibili

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// preUpload 预上传
func (uc *UploadClient) preUpload(fileName string, fileSize int64) (*PreUploadInfo, error) {
	params := url.Values{}
	params.Set("r", "upos")
	params.Set("profile", "ugcupos/bup")
	params.Set("ssl", "0")
	params.Set("version", "2.11.0")
	params.Set("build", "2110000")
	params.Set("name", fileName)
	params.Set("size", strconv.FormatInt(fileSize, 10))

	// 添加必要的参数 - 参考 biliup-rs 的实现
	params.Set("probe_version", "20221109")
	// upcdn: bda2 表示使用百度云 CDN
	// zone: cs 表示云存储区域
	params.Set("upcdn", "bda2")
	params.Set("zone", "cs")

	// 获取cookies用于认证
	cookies := uc.loginInfo.GetCookieString()

	req, err := http.NewRequest("GET", "https://member.bilibili.com/preupload?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", cookies)
	req.Header.Set("User-Agent", uc.client.userAgent)

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

	log.Printf("PreUpload API Response: %s", string(body))

	var preUploadResp PreUploadInfo
	if err := json.Unmarshal(body, &preUploadResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if preUploadResp.OK != 1 {
		return nil, fmt.Errorf("pre-upload failed: %+v", preUploadResp)
	}

	if preUploadResp.UposUri == "" {
		return nil, fmt.Errorf("pre-upload upos_uri is empty: %+v", preUploadResp)
	}

	return &preUploadResp, nil
}

// getUploadID 获取上传ID
func (uc *UploadClient) getUploadID(preInfo *PreUploadInfo) (string, error) {
	if preInfo == nil {
		return "", fmt.Errorf("preInfo is nil")
	}

	if preInfo.UposUri == "" {
		return "", fmt.Errorf("UposUri is empty")
	}

	uploadURL := fmt.Sprintf("https:%s/%s?uploads&output=json",
		preInfo.Endpoint,
		strings.Replace(preInfo.UposUri, "upos://", "", 1))

	req, err := http.NewRequest("POST", uploadURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Upos-Auth", preInfo.Auth)
	req.Header.Set("User-Agent", uc.client.userAgent)

	resp, err := uc.client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	uploadID, ok := result["upload_id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get upload_id from response: %+v", result)
	}

	return uploadID, nil
}

// uploadChunks 分块上传
func (uc *UploadClient) uploadChunks(videoPath string, preInfo *PreUploadInfo, uploadID string) ([]map[string]interface{}, error) {
	file, err := os.Open(videoPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	chunkSize := int64(preInfo.ChunkSize)
	chunksNum := int((fileSize + chunkSize - 1) / chunkSize) // 向上取整

	log.Printf("Starting chunk upload: fileSize=%d, chunkSize=%d, chunksNum=%d", fileSize, chunkSize, chunksNum)

	uploadURL := fmt.Sprintf("https:%s/%s",
		preInfo.Endpoint,
		strings.Replace(preInfo.UposUri, "upos://", "", 1))

	var parts []map[string]interface{}

	for i := 0; i < chunksNum; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if end > fileSize {
			end = fileSize
		}

		chunkData := make([]byte, end-start)
		_, err := file.ReadAt(chunkData, start)
		if err != nil {
			return nil, fmt.Errorf("failed to read chunk %d: %v", i, err)
		}

		// 上传分块（带重试）
		params := url.Values{}
		params.Set("uploadId", uploadID)
		params.Set("chunks", strconv.Itoa(chunksNum))
		params.Set("total", strconv.FormatInt(fileSize, 10))
		params.Set("chunk", strconv.Itoa(i))
		params.Set("size", strconv.Itoa(len(chunkData)))
		params.Set("partNumber", strconv.Itoa(i+1))
		params.Set("start", strconv.FormatInt(start, 10))
		params.Set("end", strconv.FormatInt(end, 10))

		// 计算上传进度
		progress := float64(i+1) / float64(chunksNum) * 100
		log.Printf("📤 Uploading chunk %d/%d (%.1f%%) - bytes %d-%d, size=%d", 
			i+1, chunksNum, progress, start, end, len(chunkData))

		err = retryFunc(func() error {
			req, err := http.NewRequest("PUT", uploadURL+"?"+params.Encode(), bytes.NewReader(chunkData))
			if err != nil {
				return fmt.Errorf("failed to create request: %v", err)
			}

			req.Header.Set("X-Upos-Auth", preInfo.Auth)
			req.Header.Set("Content-Length", strconv.Itoa(len(chunkData)))
			req.Header.Set("User-Agent", uc.client.userAgent)
			req.Header.Set("Connection", "keep-alive")

			// 使用专门的上传客户端（更长的超时时间）
			resp, err := uc.uploadClient.Do(req)
			if err != nil {
				return fmt.Errorf("network error uploading chunk %d: %v", i+1, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("upload chunk %d failed with status %s: %s", i+1, resp.Status, string(bodyBytes))
			}

			log.Printf("✅ Chunk %d/%d uploaded successfully (%.1f%% complete)", i+1, chunksNum, progress)
			return nil
		})

		if err != nil {
			return nil, err
		}

		parts = append(parts, map[string]interface{}{
			"partNumber": i + 1,
			"eTag":       "etag",
		})
	}

	log.Printf("All %d chunks uploaded successfully", chunksNum)
	return parts, nil
}

// uploadChunksFromURL 从URL流式分块上传文件到 Bilibili
func (uc *UploadClient) uploadChunksFromURL(videoURL string, preInfo *PreUploadInfo, uploadID string, fileSize int64) ([]map[string]interface{}, error) {
	chunkSize := int64(preInfo.ChunkSize)
	chunksNum := int((fileSize + chunkSize - 1) / chunkSize) // 向上取整

	log.Printf("Starting chunk upload from URL: url=%s, fileSize=%d, chunkSize=%d, chunksNum=%d", videoURL, fileSize, chunkSize, chunksNum)

	uploadURL := fmt.Sprintf("https:%s/%s",
		preInfo.Endpoint,
		strings.Replace(preInfo.UposUri, "upos://", "", 1))

	var parts []map[string]interface{}

	for i := 0; i < chunksNum; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if end > fileSize {
			end = fileSize
		}

		// 从URL下载指定范围的数据块
		chunkData, err := uc.downloadChunkFromURL(videoURL, start, end-1)
		if err != nil {
			return nil, fmt.Errorf("failed to download chunk %d: %v", i, err)
		}

		// 上传分块（带重试）
		params := url.Values{}
		params.Set("uploadId", uploadID)
		params.Set("chunks", strconv.Itoa(chunksNum))
		params.Set("total", strconv.FormatInt(fileSize, 10))
		params.Set("chunk", strconv.Itoa(i))
		params.Set("size", strconv.Itoa(len(chunkData)))
		params.Set("partNumber", strconv.Itoa(i+1))
		params.Set("start", strconv.FormatInt(start, 10))
		params.Set("end", strconv.FormatInt(end, 10))

		// 计算上传进度
		progress := float64(i+1) / float64(chunksNum) * 100
		log.Printf("📤 Uploading chunk from URL %d/%d (%.1f%%) - bytes %d-%d, size=%d", 
			i+1, chunksNum, progress, start, end, len(chunkData))

		err = retryFunc(func() error {
			req, err := http.NewRequest("PUT", uploadURL+"?"+params.Encode(), bytes.NewReader(chunkData))
			if err != nil {
				return fmt.Errorf("failed to create request: %v", err)
			}

			req.Header.Set("X-Upos-Auth", preInfo.Auth)
			req.Header.Set("Content-Length", strconv.Itoa(len(chunkData)))
			req.Header.Set("User-Agent", uc.client.userAgent)
			req.Header.Set("Connection", "keep-alive")

			// 使用专门的上传客户端（更长的超时时间）
			resp, err := uc.uploadClient.Do(req)
			if err != nil {
				return fmt.Errorf("network error uploading chunk %d: %v", i+1, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("upload chunk %d failed with status %s: %s", i+1, resp.Status, string(bodyBytes))
			}

			log.Printf("✅ Chunk from URL %d/%d uploaded successfully (%.1f%% complete)", i+1, chunksNum, progress)
			return nil
		})

		if err != nil {
			return nil, err
		}

		parts = append(parts, map[string]interface{}{
			"partNumber": i + 1,
			"eTag":       "etag",
		})
	}

	log.Printf("All %d chunks uploaded from URL successfully", chunksNum)
	return parts, nil
}

// chunkUploadTask 分块上传任务
type chunkUploadTask struct {
	index     int
	start     int64
	end       int64
	chunkData []byte
}

// chunkUploadResult 分块上传结果
type chunkUploadResult struct {
	index       int
	partNumber  int
	success     bool
	err         error
	retryCount  int
}

// endpointHealthCheck endpoint健康检查结果
type endpointHealthCheck struct {
	endpoint string
	latency  time.Duration
	available bool
}

// selectBestEndpoint 选择最优的上传endpoint
func (uc *UploadClient) selectBestEndpoint(endpoints []string) string {
	if len(endpoints) == 0 {
		return ""
	}
	
	if len(endpoints) == 1 {
		return endpoints[0]
	}
	
	// 并发检查所有endpoint的延迟
	results := make(chan endpointHealthCheck, len(endpoints))
	var wg sync.WaitGroup
	
	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(ep string) {
			defer wg.Done()
			
			// 构造健康检查URL
			probeURL := fmt.Sprintf("https:%s/OK", ep)
			
			start := time.Now()
			req, err := http.NewRequest("GET", probeURL, nil)
			if err != nil {
				results <- endpointHealthCheck{endpoint: ep, available: false}
				return
			}
			
			// 设置短超时时间进行探测
			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Do(req)
			latency := time.Since(start)
			
			if err != nil || resp.StatusCode != http.StatusOK {
				results <- endpointHealthCheck{endpoint: ep, available: false}
				return
			}
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			
			results <- endpointHealthCheck{
				endpoint:  ep,
				latency:   latency,
				available: true,
			}
		}(endpoint)
	}
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// 收集结果并选择延迟最低的可用endpoint
	var checks []endpointHealthCheck
	for check := range results {
		if check.available {
			checks = append(checks, check)
		}
	}
	
	if len(checks) == 0 {
		log.Printf("⚠️ No available endpoints, using first one: %s", endpoints[0])
		return endpoints[0]
	}
	
	// 按延迟排序
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].latency < checks[j].latency
	})
	
	bestEndpoint := checks[0].endpoint
	log.Printf("✅ Selected best endpoint: %s (latency: %v)", bestEndpoint, checks[0].latency)
	return bestEndpoint
}

// downloadChunkFromURL 从URL下载指定范围的数据块
func (uc *UploadClient) downloadChunkFromURL(url string, start, end int64) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置 Range 头
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := uc.uploadClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download chunk: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return nil, fmt.Errorf("unexpected status code: %d, expected 206", resp.StatusCode)
	}

	chunk, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read chunk data: %v", err)
	}

	return chunk, nil
}

// completeUpload 完成上传
func (uc *UploadClient) completeUpload(preInfo *PreUploadInfo, uploadID string, parts []map[string]interface{}, fileName string) (*Video, error) {
	uploadURL := fmt.Sprintf("https:%s/%s",
		preInfo.Endpoint,
		strings.Replace(preInfo.UposUri, "upos://", "", 1))

	params := url.Values{}
	params.Set("name", fileName)
	params.Set("uploadId", uploadID)
	params.Set("biz_id", strconv.FormatInt(preInfo.BizId, 10))
	params.Set("output", "json")
	params.Set("profile", "ugcupos/bup")

	requestBody := map[string]interface{}{
		"parts": parts,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Complete upload URL: %s?%s", uploadURL, params.Encode())

	req, err := http.NewRequest("POST", uploadURL+"?"+params.Encode(), bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Upos-Auth", preInfo.Auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", uc.client.userAgent)

	resp, err := uc.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Complete upload response (status %d): %s", resp.StatusCode, string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result["OK"] != float64(1) {
		return nil, fmt.Errorf("complete upload failed: %+v", result)
	}

	// 从upos_uri提取文件名（不包含扩展名）
	// upos_uri格式: upos://xxx/xxx/filename.mp4
	// 需要提取 filename 部分（不包含.mp4）
	uposPath := preInfo.UposUri
	// 移除 upos:// 前缀
	uposPath = strings.TrimPrefix(uposPath, "upos://")
	// 获取文件名部分
	baseName := filepath.Base(uposPath)
	// 移除扩展名
	fileNameWithoutExt := baseName
	if ext := filepath.Ext(baseName); ext != "" {
		fileNameWithoutExt = strings.TrimSuffix(baseName, ext)
	}

	log.Printf("Extracted filename from upos_uri '%s': '%s'", preInfo.UposUri, fileNameWithoutExt)

	return &Video{
		Title:    "", // 将在 UploadVideo 中设置
		Filename: fileNameWithoutExt,
		Desc:     "",
		CID:      preInfo.BizId,
	}, nil
}

// uploadChunksFromURLConcurrent 并发版本：从URL流式分块上传文件到 Bilibili
func (uc *UploadClient) uploadChunksFromURLConcurrent(videoURL string, preInfo *PreUploadInfo, uploadID string, fileSize int64, concurrency int) ([]map[string]interface{}, error) {
	if concurrency <= 0 {
		concurrency = 3 // 默认3个并发
	}
	
	chunkSize := int64(preInfo.ChunkSize)
	chunksNum := int((fileSize + chunkSize - 1) / chunkSize) // 向上取整

	log.Printf("🚀 Starting concurrent chunk upload from URL: url=%s, fileSize=%d, chunkSize=%d, chunksNum=%d, concurrency=%d", 
		videoURL, fileSize, chunkSize, chunksNum, concurrency)

	// 选择最优endpoint
	var uploadURL string
	if len(preInfo.Endpoints) > 1 {
		bestEndpoint := uc.selectBestEndpoint(preInfo.Endpoints)
		uploadURL = fmt.Sprintf("https:%s/%s", bestEndpoint, strings.Replace(preInfo.UposUri, "upos://", "", 1))
	} else {
		uploadURL = fmt.Sprintf("https:%s/%s", preInfo.Endpoint, strings.Replace(preInfo.UposUri, "upos://", "", 1))
	}

	// 创建任务通道和结果通道
	tasks := make(chan chunkUploadTask, chunksNum)
	results := make(chan chunkUploadResult, chunksNum)
	
	// 用于跟踪整体进度
	var progressMutex sync.Mutex
	completedChunks := 0
	totalBytes := int64(0)
	startTime := time.Now()
	
	// 启动worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for task := range tasks {
				// 从URL下载指定范围的数据块
				chunkData, err := uc.downloadChunkFromURL(videoURL, task.start, task.end-1)
				if err != nil {
					results <- chunkUploadResult{
						index:   task.index,
						success: false,
						err:     fmt.Errorf("failed to download chunk %d: %v", task.index, err),
					}
					continue
				}
				
				// 上传分块（带重试）
				params := url.Values{}
				params.Set("uploadId", uploadID)
				params.Set("chunks", strconv.Itoa(chunksNum))
				params.Set("total", strconv.FormatInt(fileSize, 10))
				params.Set("chunk", strconv.Itoa(task.index))
				params.Set("size", strconv.Itoa(len(chunkData)))
				params.Set("partNumber", strconv.Itoa(task.index+1))
				params.Set("start", strconv.FormatInt(task.start, 10))
				params.Set("end", strconv.FormatInt(task.end, 10))
				
				retryCount := 0
				err = retryFunc(func() error {
					retryCount++
					req, err := http.NewRequest("PUT", uploadURL+"?"+params.Encode(), bytes.NewReader(chunkData))
					if err != nil {
						return fmt.Errorf("failed to create request: %v", err)
					}

					req.Header.Set("X-Upos-Auth", preInfo.Auth)
					req.Header.Set("Content-Length", strconv.Itoa(len(chunkData)))
					req.Header.Set("User-Agent", uc.client.userAgent)
					req.Header.Set("Connection", "keep-alive")

					// 使用专门的上传客户端
					resp, err := uc.uploadClient.Do(req)
					if err != nil {
						return fmt.Errorf("network error uploading chunk %d: %v", task.index+1, err)
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						bodyBytes, _ := io.ReadAll(resp.Body)
						return fmt.Errorf("upload chunk %d failed with status %s: %s", task.index+1, resp.Status, string(bodyBytes))
					}
					
					return nil
				})
				
				// 更新进度
				progressMutex.Lock()
				completedChunks++
				totalBytes += int64(len(chunkData))
				progress := float64(completedChunks) / float64(chunksNum) * 100
				elapsed := time.Since(startTime)
				speed := float64(totalBytes) / elapsed.Seconds() / 1024 / 1024 // MB/s
				eta := time.Duration(float64(fileSize-totalBytes)/float64(totalBytes)*elapsed.Seconds()) * time.Second
				
				if err == nil {
					log.Printf("✅ [Worker-%d] Chunk %d/%d uploaded (%.1f%%) | Speed: %.2f MB/s | ETA: %v", 
						workerID, completedChunks, chunksNum, progress, speed, eta.Round(time.Second))
				} else {
					log.Printf("❌ [Worker-%d] Chunk %d/%d failed after %d retries: %v", 
						workerID, task.index+1, chunksNum, retryCount, err)
				}
				progressMutex.Unlock()
				
				results <- chunkUploadResult{
					index:      task.index,
					partNumber: task.index + 1,
					success:    err == nil,
					err:        err,
					retryCount: retryCount,
				}
			}
		}(i)
	}
	
	// 生成所有上传任务
	go func() {
		for i := 0; i < chunksNum; i++ {
			start := int64(i) * chunkSize
			end := start + chunkSize
			if end > fileSize {
				end = fileSize
			}
			
			tasks <- chunkUploadTask{
				index: i,
				start: start,
				end:   end,
			}
		}
		close(tasks)
	}()
	
	// 等待所有worker完成
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// 收集结果
	var parts []map[string]interface{}
	var errors []error
	resultMap := make(map[int]chunkUploadResult)
	
	for result := range results {
		resultMap[result.index] = result
		if !result.success {
			errors = append(errors, result.err)
		}
	}
	
	// 检查是否有失败
	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to upload %d chunks: %v", len(errors), errors[0])
	}
	
	// 按索引顺序构建parts数组
	for i := 0; i < chunksNum; i++ {
		parts = append(parts, map[string]interface{}{
			"partNumber": i + 1,
			"eTag":       "etag",
		})
	}
	
	totalTime := time.Since(startTime)
	avgSpeed := float64(fileSize) / totalTime.Seconds() / 1024 / 1024
	log.Printf("🎉 All %d chunks uploaded successfully! Total time: %v, Average speed: %.2f MB/s", 
		chunksNum, totalTime.Round(time.Second), avgSpeed)
	
	return parts, nil
}