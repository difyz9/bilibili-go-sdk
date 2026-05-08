package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

func main() {
	// 1. 从环境变量或配置文件加载登录信息
	loginInfo := loadLoginInfo()

	// 2. 创建上传客户端
	uploadClient := bilibili.NewUploadClient(loginInfo)

	// 3. 腾讯COS视频URL和信息
	cosURL := "https://your-bucket.cos.ap-guangzhou.myqcloud.com/video.mp4"
	fileName := "my-video.mp4"
	fileSize := int64(305333744) // 文件大小，需提前获取

	// 4. 使用并发上传（推荐）
	// 参数说明：
	// - cosURL: 腾讯COS的视频链接
	// - fileName: 文件名
	// - fileSize: 文件大小（字节）
	// - concurrency: 并发数，建议3-5，设为0使用默认值3
	log.Println("开始并发上传视频...")
	video, err := uploadClient.UploadVideoFromURLConcurrent(cosURL, fileName, fileSize, 3)
	if err != nil {
		log.Fatalf("上传失败: %v", err)
	}

	log.Printf("✅ 视频上传成功！文件名: %s", video.Filename)

	// 5. 提交视频信息到B站
	studio := &bilibili.Studio{
		Title:         "我的视频标题",
		Tid:           122, // 分区ID，122=野生技能协会
		Tag:           "标签1,标签2,标签3",
		Desc:          "这是视频描述",
		Copyright:     1, // 1=自制，2=转载
		Source:        "",
		Cover:         "https://your-cover-url.jpg", // 封面图URL
		Videos:        []bilibili.Video{*video},
		NoReprint:     1,  // 1=允许转载，0=不允许
		OpenElec:      1,  // 1=开启充电，0=关闭
		Recreate:      -1, // -1=允许二创，1=不允许
		Dynamic:       "动态内容",
		Interactive:   0,
		NoDisturbance: 1, // 1=推送到动态，0=不推送
	}

	result, err := uploadClient.SubmitVideo(studio)
	if err != nil {
		log.Fatalf("提交视频失败: %v", err)
	}

	log.Printf("🎉 视频发布成功！")
	log.Printf("   稿件ID: %v", result.Data)
}

func loadLoginInfo() *bilibili.LoginInfo {
	// 这里应该从配置文件或环境变量加载实际的登录信息
	// 示例结构：
	return &bilibili.LoginInfo{
		CookieInfo: map[string]interface{}{
			"cookies": []interface{}{
				map[string]interface{}{
					"name":  "SESSDATA",
					"value": os.Getenv("BILI_SESSDATA"),
				},
				map[string]interface{}{
					"name":  "bili_jct",
					"value": os.Getenv("BILI_JCT"),
				},
				map[string]interface{}{
					"name":  "DedeUserID",
					"value": os.Getenv("BILI_UID"),
				},
			},
		},
	}
}

// 性能对比示例
func performanceComparison() {
	loginInfo := loadLoginInfo()
	uploadClient := bilibili.NewUploadClient(loginInfo)

	cosURL := "https://your-bucket.cos.ap-guangzhou.myqcloud.com/video.mp4"
	fileName := "test-video.mp4"
	fileSize := int64(305333744)

	// 方法1：单线程上传（原版）
	fmt.Println("=== 单线程上传测试 ===")
	video1, err := uploadClient.UploadVideoFromURL(cosURL, fileName, fileSize)
	if err != nil {
		log.Printf("单线程上传失败: %v", err)
	} else {
		log.Printf("单线程上传成功: %s", video1.Filename)
	}

	// 方法2：3并发上传（推荐）
	fmt.Println("\n=== 3并发上传测试 ===")
	video2, err := uploadClient.UploadVideoFromURLConcurrent(cosURL, fileName, fileSize, 3)
	if err != nil {
		log.Printf("3并发上传失败: %v", err)
	} else {
		log.Printf("3并发上传成功: %s", video2.Filename)
	}

	// 方法3：5并发上传（网络条件好时）
	fmt.Println("\n=== 5并发上传测试 ===")
	video3, err := uploadClient.UploadVideoFromURLConcurrent(cosURL, fileName, fileSize, 5)
	if err != nil {
		log.Printf("5并发上传失败: %v", err)
	} else {
		log.Printf("5并发上传成功: %s", video3.Filename)
	}
}

// 获取腾讯COS文件大小的辅助函数
func getFileSize(cosURL string) (int64, error) {
	// 方法1：HEAD请求获取Content-Length
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", cosURL, nil)
	if err != nil {
		return 0, err
	}

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
