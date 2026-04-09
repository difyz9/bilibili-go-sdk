package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

func main() {
	// 命令行参数
	videoURL := flag.String("url", "", "视频URL（必需）")
	fileName := flag.String("name", "test-video.mp4", "文件名")
	mode := flag.String("mode", "concurrent", "上传模式: single, concurrent, compare")
	concurrency := flag.Int("c", 3, "并发数（仅concurrent模式）")
	flag.Parse()

	if *videoURL == "" {
		log.Fatal("请提供视频URL: -url https://your-video-url.mp4")
	}

	// 获取登录信息
	loginInfo := getLoginInfo()

	// 获取文件大小
	fileSize, err := getFileSize(*videoURL)
	if err != nil {
		log.Fatalf("获取文件大小失败: %v", err)
	}

	log.Printf("📁 视频信息:")
	log.Printf("   URL: %s", *videoURL)
	log.Printf("   文件名: %s", *fileName)
	log.Printf("   大小: %.2f MB", float64(fileSize)/1024/1024)
	log.Printf("   模式: %s", *mode)
	if *mode == "concurrent" {
		log.Printf("   并发数: %d", *concurrency)
	}
	log.Println()

	// 创建上传客户端
	uploadClient := bilibili.NewUploadClient(loginInfo)

	switch *mode {
	case "single":
		testSingleThread(uploadClient, *videoURL, *fileName, fileSize)
	case "concurrent":
		testConcurrent(uploadClient, *videoURL, *fileName, fileSize, *concurrency)
	case "compare":
		comparePerformance(uploadClient, *videoURL, *fileName, fileSize)
	default:
		log.Fatal("无效的模式，请使用: single, concurrent, compare")
	}
}

// testSingleThread 测试单线程上传
func testSingleThread(client *bilibili.UploadClient, url, name string, size int64) {
	log.Println("========================================")
	log.Println("🐢 单线程上传测试")
	log.Println("========================================")

	start := time.Now()
	video, err := client.UploadVideoFromURL(url, name, size)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("❌ 上传失败: %v", err)
		return
	}

	printResult("单线程", video, size, elapsed)
}

// testConcurrent 测试并发上传
func testConcurrent(client *bilibili.UploadClient, url, name string, size int64, concurrency int) {
	log.Println("========================================")
	log.Printf("🚀 %d并发上传测试", concurrency)
	log.Println("========================================")

	start := time.Now()
	video, err := client.UploadVideoFromURLConcurrent(url, name, size, concurrency)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("❌ 上传失败: %v", err)
		return
	}

	printResult(fmt.Sprintf("%d并发", concurrency), video, size, elapsed)
}

// comparePerformance 性能对比测试
func comparePerformance(client *bilibili.UploadClient, url, name string, size int64) {
	log.Println("========================================")
	log.Println("📊 性能对比测试")
	log.Println("========================================")

	results := make(map[string]time.Duration)

	// 测试1: 单线程
	log.Println("\n[测试 1/4] 单线程上传...")
	start1 := time.Now()
	_, err1 := client.UploadVideoFromURL(url, name, size)
	elapsed1 := time.Since(start1)
	if err1 == nil {
		results["单线程"] = elapsed1
		log.Printf("✅ 单线程完成: %v", elapsed1.Round(time.Second))
	} else {
		log.Printf("❌ 单线程失败: %v", err1)
	}

	time.Sleep(2 * time.Second) // 等待2秒

	// 测试2: 2并发
	log.Println("\n[测试 2/4] 2并发上传...")
	start2 := time.Now()
	_, err2 := client.UploadVideoFromURLConcurrent(url, name, size, 2)
	elapsed2 := time.Since(start2)
	if err2 == nil {
		results["2并发"] = elapsed2
		log.Printf("✅ 2并发完成: %v", elapsed2.Round(time.Second))
	} else {
		log.Printf("❌ 2并发失败: %v", err2)
	}

	time.Sleep(2 * time.Second)

	// 测试3: 3并发
	log.Println("\n[测试 3/4] 3并发上传...")
	start3 := time.Now()
	_, err3 := client.UploadVideoFromURLConcurrent(url, name, size, 3)
	elapsed3 := time.Since(start3)
	if err3 == nil {
		results["3并发"] = elapsed3
		log.Printf("✅ 3并发完成: %v", elapsed3.Round(time.Second))
	} else {
		log.Printf("❌ 3并发失败: %v", err3)
	}

	time.Sleep(2 * time.Second)

	// 测试4: 5并发
	log.Println("\n[测试 4/4] 5并发上传...")
	start4 := time.Now()
	_, err4 := client.UploadVideoFromURLConcurrent(url, name, size, 5)
	elapsed4 := time.Since(start4)
	if err4 == nil {
		results["5并发"] = elapsed4
		log.Printf("✅ 5并发完成: %v", elapsed4.Round(time.Second))
	} else {
		log.Printf("❌ 5并发失败: %v", err4)
	}

	// 打印对比结果
	log.Println("\n========================================")
	log.Println("📈 性能对比结果")
	log.Println("========================================")

	if baseTime, ok := results["单线程"]; ok {
		fmt.Printf("\n%-10s | %-12s | %-12s | %-8s\n", "模式", "耗时", "速度", "提升")
		fmt.Println("-----------|--------------|--------------|----------")

		for _, mode := range []string{"单线程", "2并发", "3并发", "5并发"} {
			if elapsed, ok := results[mode]; ok {
				speed := float64(size) / elapsed.Seconds() / 1024 / 1024
				improvement := float64(baseTime) / float64(elapsed)
				fmt.Printf("%-10s | %-12s | %8.2f MB/s | %.2fx\n",
					mode,
					elapsed.Round(time.Second),
					speed,
					improvement,
				)
			}
		}

		// 推荐建议
		log.Println("\n========================================")
		log.Println("💡 推荐建议")
		log.Println("========================================")

		best := "单线程"
		bestTime := baseTime
		for mode, elapsed := range results {
			if elapsed < bestTime {
				bestTime = elapsed
				best = mode
			}
		}

		log.Printf("🏆 最快模式: %s (耗时 %v)", best, bestTime.Round(time.Second))

		if best != "单线程" {
			improvement := float64(baseTime) / float64(bestTime)
			log.Printf("⚡ 性能提升: %.2fx", improvement)
		}
	}
}

// printResult 打印上传结果
func printResult(mode string, video *bilibili.Video, size int64, elapsed time.Duration) {
	log.Println("\n========================================")
	log.Println("✅ 上传成功")
	log.Println("========================================")
	log.Printf("模式: %s", mode)
	log.Printf("文件名: %s", video.Filename)
	log.Printf("标题: %s", video.Title)
	log.Printf("耗时: %v", elapsed.Round(time.Second))
	log.Printf("文件大小: %.2f MB", float64(size)/1024/1024)
	log.Printf("平均速度: %.2f MB/s", float64(size)/elapsed.Seconds()/1024/1024)
	log.Println("========================================")
}

// getFileSize 获取文件大小
func getFileSize(url string) (int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HEAD request failed: %s", resp.Status)
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("no Content-Length header")
	}

	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid Content-Length: %v", err)
	}

	return size, nil
}

// getLoginInfo 获取登录信息
func getLoginInfo() *bilibili.LoginInfo {
	sessdata := os.Getenv("BILI_SESSDATA")
	if sessdata == "" {
		log.Fatal("请设置环境变量 BILI_SESSDATA")
	}

	biliJct := os.Getenv("BILI_JCT")
	if biliJct == "" {
		log.Fatal("请设置环境变量 BILI_JCT")
	}

	return &bilibili.LoginInfo{
		CookieInfo: map[string]interface{}{
			"cookies": []interface{}{
				map[string]interface{}{
					"name":  "SESSDATA",
					"value": sessdata,
				},
				map[string]interface{}{
					"name":  "bili_jct",
					"value": biliJct,
				},
			},
		},
	}
}
