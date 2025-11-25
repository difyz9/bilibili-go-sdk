// 用户信息和粉丝数示例
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

const loginInfoFile = "login.json"

func main() {
	fmt.Println("=== 获取用户信息和粉丝数 ===\n")

	// 方式1: 使用已保存的登录信息获取当前用户信息
	loginInfo, err := loadLoginInfo()
	if err != nil {
		fmt.Printf("❌ 加载登录信息失败: %v\n", err)
		fmt.Println("请先运行登录示例: go run examples/login/main.go")
		fmt.Println("\n=== 继续演示方式2: 获取其他用户的公开信息 ===\n")
		demonstratePublicUserInfo()
		return
	}

	// 创建客户端
	client := bilibili.NewClient()

	// 获取当前登录用户的详细信息
	cookies := loginInfo.GetCookieString()
	myInfo, err := client.GetMyInfoWithRetry(cookies, 3)
	if err != nil {
		log.Fatalf("❌ 获取用户信息失败: %v", err)
	}

	// 显示用户信息
	fmt.Println("✅ 当前登录用户信息:")
	fmt.Printf("  - 用户ID (Mid): %d\n", myInfo.Mid)
	fmt.Printf("  - 用户名: %s\n", myInfo.Name)
	fmt.Printf("  - 粉丝数: %d\n", myInfo.Follower)
	fmt.Printf("  - 关注数: %d\n", myInfo.Following)
	fmt.Printf("  - 等级: %d\n", myInfo.Level)
	fmt.Printf("  - 硬币数: %d\n", myInfo.GetCoins())
	fmt.Printf("  - 签名: %s\n", myInfo.Sign)
	fmt.Printf("  - 头像: %s\n", myInfo.Face)

	// 也可以使用 UserStat API 获取关系统计
	fmt.Println("\n--- 使用 GetUserStat 获取关系统计 ---")
	userStat, err := client.GetUserStat(myInfo.Mid)
	if err != nil {
		log.Printf("获取关系统计失败: %v", err)
	} else {
		fmt.Printf("  - 粉丝数: %d\n", userStat.Follower)
		fmt.Printf("  - 关注数: %d\n", userStat.Following)
		fmt.Printf("  - 悄悄关注: %d\n", userStat.Whisper)
		fmt.Printf("  - 黑名单: %d\n", userStat.Black)
	}

	// 演示获取其他用户的公开信息
	fmt.Println("\n=== 获取其他用户的公开信息 ===\n")
	demonstratePublicUserInfo()
}

// demonstratePublicUserInfo 演示如何获取其他用户的公开信息
func demonstratePublicUserInfo() {
	client := bilibili.NewClient()

	// 示例：获取某个用户的粉丝数和统计信息
	// 这里使用一个示例用户ID，你可以替换为任何bilibili用户的mid
	exampleMid := int64(3549361111093) // 示例用户ID

	fmt.Printf("获取用户 %d 的公开统计信息...\n\n", exampleMid)

	// 获取用户关系统计（粉丝数、关注数等）
	userStat, err := client.GetUserStat(exampleMid)
	if err != nil {
		log.Printf("❌ 获取用户统计失败: %v", err)
		return
	}

	fmt.Println("✅ 用户关系统计:")
	fmt.Printf("  - 用户ID: %d\n", userStat.Mid)
	fmt.Printf("  - 粉丝数: %d\n", userStat.Follower)
	fmt.Printf("  - 关注数: %d\n", userStat.Following)

	// 注意：GetUserStat 只能获取统计数据，不包含用户名
	// 如果需要获取用户名，需要调用其他API（如用户空间API）
	// SDK中可以通过访问用户空间来获取更多信息

	fmt.Println("\n💡 提示:")
	fmt.Println("  - 使用 GetMyInfo() 可以获取当前登录用户的完整信息（包括用户名、粉丝数）")
	fmt.Println("  - 使用 GetUserStat() 可以获取任何用户的公开统计信息（粉丝数、关注数）")
	fmt.Println("  - 获取其他用户的用户名需要调用用户空间相关API")
}

// loadLoginInfo 从文件加载登录信息
func loadLoginInfo() (*bilibili.LoginInfo, error) {
	data, err := ioutil.ReadFile(loginInfoFile)
	if err != nil {
		return nil, err
	}

	// 先尝试解析包含 login_info 字段的结构
	var wrapper struct {
		LoginInfo *bilibili.LoginInfo `json:"login_info"`
	}
	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.LoginInfo != nil {
		return wrapper.LoginInfo, nil
	}

	// 如果失败，尝试直接解析为 LoginInfo
	var loginInfo bilibili.LoginInfo
	if err := json.Unmarshal(data, &loginInfo); err != nil {
		return nil, err
	}

	return &loginInfo, nil
}
