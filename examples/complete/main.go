// 完整流程示例 - 登录、上传、投稿
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

const loginInfoFile = "login_info.json"


// go run complete_example.go upload <视频路径> [封面路径] # 上传视频

// go run examples/complete/main.go upload examples/001.mp4 examples/002.jpg

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("用法:")
		fmt.Println("  go run examples/complete/main.go login                    # 登录")
		fmt.Println("  go run examples/complete/main.go upload examples/001.mp4 examples/002.jpg # 上传视频")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "login":
		doLogin()
	case "upload":
		if len(os.Args) < 3 {
			fmt.Println("请提供视频文件路径")
			os.Exit(1)
		}
		videoPath := os.Args[2]
		var coverPath string
		if len(os.Args) >= 4 {
			coverPath = os.Args[3]
		}
		doUpload(videoPath, coverPath)
	default:
		fmt.Printf("未知命令: %s\n", command)
		os.Exit(1)
	}
}

func doLogin() {
	fmt.Println("=== Bilibili 登录 ===")

	client := bilibili.NewClient()

	// 获取二维码
	qrResp, err := client.GetQRCode()
	if err != nil {
		log.Fatalf("获取二维码失败: %v", err)
	}

	fmt.Printf("请用手机bilibili扫描二维码登录:\n%s\n", qrResp.Data.URL)

	// 轮询登录状态
	loginInfo, err := client.PollQRCode(qrResp.Data.AuthCode)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("登录成功! 用户: %s (ID: %d)\n", loginInfo.TokenInfo.Uname, loginInfo.TokenInfo.Mid)

	// 保存登录信息到文件
	if err := saveLoginInfo(loginInfo); err != nil {
		log.Printf("保存登录信息失败: %v", err)
	} else {
		fmt.Printf("登录信息已保存到 %s\n", loginInfoFile)
	}

	// 获取用户详细信息
	cookies := loginInfo.GetCookieString()
	myInfo, err := client.GetMyInfoWithRetry(cookies, 3)
	if err != nil {
		log.Printf("获取用户详细信息失败: %v", err)
	} else {
		fmt.Printf("用户详情: 等级%d, 粉丝%d, 硬币%d\n", 
			myInfo.Level, myInfo.Fans, myInfo.Coins)
	}
}

func doUpload(videoPath, coverPath string) {
	fmt.Println("=== 视频上传 ===")

	// 加载登录信息
	loginInfo, err := loadLoginInfo()
	if err != nil {
		log.Fatalf("加载登录信息失败: %v\n请先运行 'go run complete_example.go login' 登录", err)
	}

	fmt.Printf("使用已保存的登录信息 (用户: %s)\n", loginInfo.TokenInfo.Uname)

	// 创建上传客户端
	uploader := bilibili.NewUploadClient(loginInfo)

	// 上传视频
	fmt.Printf("开始上传视频: %s\n", videoPath)
	video, err := uploader.UploadVideo(videoPath)
	if err != nil {
		log.Fatalf("视频上传失败: %v", err)
	}

	fmt.Printf("✅ 视频上传完成: %s\n", video.Filename)

	// 上传封面
	var coverURL string
	if coverPath != "" {
		fmt.Printf("开始上传封面: %s\n", coverPath)
		coverURL, err = uploader.UploadCover(coverPath)
		if err != nil {
			log.Printf("❌ 封面上传失败: %v", err)
		} else {
			fmt.Printf("✅ 封面上传完成: %s\n", coverURL)
		}
	}

	// 构建投稿信息
	studio := &bilibili.Studio{
		Title:    "Bilibili Go SDK 上传测试视频",
		Desc:      "使用 Bilibili Go SDK 上传的视频\n\n这是一个测试视频，演示 SDK 的功能。",
		Tid:       122, // 生活区
		Cover:     coverURL,
		Tag:       "SDK,测试,上传,bilibili",
		Copyright: 1, // 原创
		Videos:    []bilibili.Video{*video},
		
		OpenSubtitle:  false,
		Interactive:   0,
		NoReprint:     1,
		OpenElec:      1,
		Dolby:         0,
		LosslessMusic: 0,
	}

	// 提交投稿
	fmt.Println("开始提交投稿...")
	result, err := uploader.SubmitVideo(studio)
	if err != nil {
		log.Fatalf("投稿提交失败: %v", err)
	}

	if result.Code == 0 {
		fmt.Println("🎉 投稿提交成功!")
		if data, ok := result.Data.(map[string]interface{}); ok {
			if aid, ok := data["aid"]; ok {
				fmt.Printf("视频AV号: %v\n", aid)
			}
		}
	} else {
		fmt.Printf("❌ 投稿提交失败: code=%d, message=%s\n", result.Code, result.Message)
	}
}

func saveLoginInfo(loginInfo *bilibili.LoginInfo) error {
	data, err := json.MarshalIndent(loginInfo, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(loginInfoFile, data, 0600)
}

func loadLoginInfo() (*bilibili.LoginInfo, error) {
	data, err := ioutil.ReadFile(loginInfoFile)
	if err != nil {
		return nil, err
	}

	var loginInfo bilibili.LoginInfo
	if err := json.Unmarshal(data, &loginInfo); err != nil {
		return nil, err
	}

	return &loginInfo, nil
}