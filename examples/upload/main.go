// 视频上传示例
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("用法: go run upload_example.go <视频文件路径> [封面图片路径]")
		os.Exit(1)
	}

	videoPath := os.Args[1]
	var coverPath string
	if len(os.Args) >= 3 {
		coverPath = os.Args[2]
	}

	// 这里需要先登录获取 loginInfo
	// 为了演示，我们假设已经有了 loginInfo
	// 实际使用时，需要先运行登录流程
	
	fmt.Println("注意: 此示例需要先完成登录流程获取 LoginInfo")
	fmt.Println("请参考 login_example.go 获取登录信息")
	
	// 示例 LoginInfo (实际使用时应该从登录流程获取)
	var loginInfo *bilibili.LoginInfo
	if loginInfo == nil {
		log.Fatal("请先完成登录流程获取 LoginInfo")
	}

	// 创建上传客户端
	uploader := bilibili.NewUploadClient(loginInfo)

	// 上传视频
	fmt.Printf("开始上传视频: %s\n", videoPath)
	video, err := uploader.UploadVideo(videoPath)
	if err != nil {
		log.Fatalf("视频上传失败: %v", err)
	}

	fmt.Printf("视频上传成功!\n")
	fmt.Printf("文件名: %s\n", video.Filename)
	fmt.Printf("标题: %s\n", video.Title)

	// 上传封面（如果提供了）
	var coverURL string
	if coverPath != "" {
		fmt.Printf("开始上传封面: %s\n", coverPath)
		coverURL, err = uploader.UploadCover(coverPath)
		if err != nil {
			log.Printf("封面上传失败: %v", err)
		} else {
			fmt.Printf("封面上传成功: %s\n", coverURL)
		}
	}

	// 构建投稿信息
	studio := &bilibili.Studio{
		Title:     video.Title,                    // 使用视频文件名作为标题
		Desc:      "通过 Bilibili Go SDK 上传",    // 视频描述
		Tid:       174,                           // 分区ID (174 = 生活区)
		Cover:     coverURL,                      // 封面URL
		Tag:       "测试,SDK,上传",                // 标签
		Copyright: 1,                             // 1=原创, 2=转载
		Videos:    []bilibili.Video{*video},     // 视频列表
		
		// 可选设置
		OpenSubtitle:  false,                     // 关闭字幕
		Interactive:   0,                         // 关闭互动
		NoReprint:     1,                         // 禁止转载
		OpenElec:      1,                         // 开启充电
		Dolby:         0,                         // 关闭杜比音效
		LosslessMusic: 0,                         // 关闭Hi-Res
	}

	// 提交投稿
	fmt.Println("开始提交投稿...")
	result, err := uploader.SubmitVideo(studio)
	if err != nil {
		log.Fatalf("投稿提交失败: %v", err)
	}

	if result.Code == 0 {
		fmt.Println("投稿提交成功!")
		fmt.Printf("响应: %+v\n", result)
	} else {
		fmt.Printf("投稿提交失败: code=%d, message=%s\n", result.Code, result.Message)
	}
}