// 基本登录示例
package main

import (
	"fmt"
	"log"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

func main() {
	// 创建客户端
	client := bilibili.NewClient()

	// 获取二维码
	fmt.Println("正在获取登录二维码...")
	qrResp, err := client.GetQRCode()
	if err != nil {
		log.Fatalf("获取二维码失败: %v", err)
	}

	fmt.Printf("请扫描二维码登录: %s\n", qrResp.Data.URL)
	fmt.Printf("授权码: %s\n", qrResp.Data.AuthCode)

	// 轮询登录状态
	fmt.Println("等待扫码登录...")
	loginInfo, err := client.PollQRCode(qrResp.Data.AuthCode)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
	}

	fmt.Printf("登录成功!\n")
	fmt.Printf("用户ID: %d\n", loginInfo.TokenInfo.Mid)
	fmt.Printf("用户名: %s\n", loginInfo.TokenInfo.Uname)

	// 获取详细用户信息
	fmt.Println("获取用户详细信息...")
	cookies := loginInfo.GetCookieString()
	
	myInfo, err := client.GetMyInfoWithRetry(cookies, 3)
	if err != nil {
		log.Printf("获取详细信息失败，尝试获取基本信息: %v", err)
		
		// 回退到基本信息
		userInfo, err := client.GetUserBasicInfoWithRetry(loginInfo.TokenInfo.Mid, cookies, 3)
		if err != nil {
			log.Printf("获取基本信息也失败: %v", err)
		} else {
			fmt.Printf("用户基本信息:\n")
			fmt.Printf("  - 用户名: %s\n", userInfo.Name)
			fmt.Printf("  - 等级: %d\n", userInfo.Level)
			fmt.Printf("  - 签名: %s\n", userInfo.Sign)
		}
	} else {
		fmt.Printf("用户详细信息:\n")
		fmt.Printf("  - 用户名: %s\n", myInfo.Uname)
		fmt.Printf("  - 等级: %d\n", myInfo.Level)
		fmt.Printf("  - 粉丝数: %d\n", myInfo.Fans)
		fmt.Printf("  - 关注数: %d\n", myInfo.Attention)
		fmt.Printf("  - 硬币数: %d\n", myInfo.Coins)
		fmt.Printf("  - 签名: %s\n", myInfo.Sign)
	}

	// 获取分区信息
	fmt.Println("获取投稿分区信息...")
	archiveData, err := client.GetArchivePre(cookies)
	if err != nil {
		log.Printf("获取分区信息失败: %v", err)
	} else {
		fmt.Printf("可用分区数量: %d\n", len(archiveData.TypeList))
		for i, partition := range archiveData.TypeList {
			if i < 5 { // 只显示前5个分区
				fmt.Printf("  - %s (ID: %d)\n", partition.Name, partition.ID)
			}
		}
		if len(archiveData.TypeList) > 5 {
			fmt.Printf("  ... 还有 %d 个分区\n", len(archiveData.TypeList)-5)
		}
	}

	fmt.Println("登录演示完成!")
}