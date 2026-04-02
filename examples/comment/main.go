package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

const loginInfoFile = "login.json"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cookies, err := loadCookies()
	if err != nil {
		fmt.Printf("加载 Cookies 失败: %v\n", err)
		fmt.Println("可通过环境变量 BILIBILI_COOKIES 提供，或先运行登录示例生成 login.json")
		os.Exit(1)
	}

	client := bilibili.NewClient()
	command := os.Args[1]

	switch command {
	case "add":
		handleAdd(client, cookies)
	case "like":
		handleAction(client, cookies, "like")
	case "hate":
		handleAction(client, cookies, "hate")
	case "delete":
		handleDelete(client, cookies)
	case "top":
		handleAction(client, cookies, "top")
	case "report":
		handleReport(client, cookies)
	default:
		fmt.Printf("未知命令: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleAdd(client *bilibili.Client, cookies string) {
	if len(os.Args) < 5 {
		fmt.Println("参数不足: add <type> <oid> <message> [root] [parent]")
		os.Exit(1)
	}

	commentType := mustParseInt(os.Args[2], "type")
	oid := mustParseInt64(os.Args[3], "oid")
	message := os.Args[4]

	request := &bilibili.CommentAddRequest{
		Type:    commentType,
		OID:     oid,
		Message: message,
	}
	if len(os.Args) >= 6 {
		request.Root = mustParseInt64(os.Args[5], "root")
	}
	if len(os.Args) >= 7 {
		request.Parent = mustParseInt64(os.Args[6], "parent")
	}

	response, err := client.AddComment(request, cookies)
	if err != nil {
		fmt.Printf("发表评论失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("发表评论成功")
	fmt.Printf("  rpid: %d\n", response.RPID)
	fmt.Printf("  提示: %s\n", response.SuccessToast)
	fmt.Printf("  root: %d\n", response.Root)
	fmt.Printf("  parent: %d\n", response.Parent)
}

func handleAction(client *bilibili.Client, cookies string, kind string) {
	if len(os.Args) < 5 {
		fmt.Printf("参数不足: %s <type> <oid> <rpid> [action]\n", kind)
		os.Exit(1)
	}

	request := &bilibili.CommentActionRequest{
		Type:   mustParseInt(os.Args[2], "type"),
		OID:    mustParseInt64(os.Args[3], "oid"),
		RPID:   mustParseInt64(os.Args[4], "rpid"),
		Action: 1,
	}
	if len(os.Args) >= 6 {
		request.Action = mustParseInt(os.Args[5], "action")
	}

	var err error
	switch kind {
	case "like":
		err = client.LikeComment(request, cookies)
	case "hate":
		err = client.HateComment(request, cookies)
	case "top":
		err = client.TopComment(request, cookies)
	default:
		fmt.Printf("未知操作: %s\n", kind)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("%s 评论失败: %v\n", kind, err)
		os.Exit(1)
	}

	verb := map[string]string{"like": "点赞", "hate": "点踩", "top": "置顶"}[kind]
	if request.Action == 0 {
		fmt.Printf("取消%s成功\n", verb)
	} else {
		fmt.Printf("%s成功\n", verb)
	}
	fmt.Printf("  rpid: %d\n", request.RPID)
}

func handleDelete(client *bilibili.Client, cookies string) {
	if len(os.Args) < 5 {
		fmt.Println("参数不足: delete <type> <oid> <rpid>")
		os.Exit(1)
	}

	err := client.DeleteComment(&bilibili.CommentDeleteRequest{
		Type: mustParseInt(os.Args[2], "type"),
		OID:  mustParseInt64(os.Args[3], "oid"),
		RPID: mustParseInt64(os.Args[4], "rpid"),
	}, cookies)
	if err != nil {
		fmt.Printf("删除评论失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("删除评论成功")
}

func handleReport(client *bilibili.Client, cookies string) {
	if len(os.Args) < 6 {
		fmt.Println("参数不足: report <type> <oid> <rpid> <reason> [content]")
		os.Exit(1)
	}

	content := ""
	if len(os.Args) >= 7 {
		content = os.Args[6]
	}

	err := client.ReportComment(&bilibili.CommentReportRequest{
		Type:    mustParseInt(os.Args[2], "type"),
		OID:     mustParseInt64(os.Args[3], "oid"),
		RPID:    mustParseInt64(os.Args[4], "rpid"),
		Reason:  mustParseInt(os.Args[5], "reason"),
		Content: content,
	}, cookies)
	if err != nil {
		fmt.Printf("举报评论失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("举报评论成功")
}

func loadCookies() (string, error) {
	if cookies := strings.TrimSpace(os.Getenv("BILIBILI_COOKIES")); cookies != "" {
		return cookies, nil
	}

	data, err := os.ReadFile(loginInfoFile)
	if err != nil {
		return "", err
	}

	var wrapper struct {
		LoginInfo *bilibili.LoginInfo `json:"login_info"`
	}
	if err := json.Unmarshal(data, &wrapper); err == nil && wrapper.LoginInfo != nil {
		return wrapper.LoginInfo.GetCookieString(), nil
	}

	var loginInfo bilibili.LoginInfo
	if err := json.Unmarshal(data, &loginInfo); err != nil {
		return "", err
	}

	return loginInfo.GetCookieString(), nil
}

func mustParseInt(value string, field string) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("无效的 %s: %s\n", field, value)
		os.Exit(1)
	}
	return parsed
}

func mustParseInt64(value string, field string) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Printf("无效的 %s: %s\n", field, value)
		os.Exit(1)
	}
	return parsed
}

func printUsage() {
	fmt.Println("评论操作示例")
	fmt.Println("")
	fmt.Println("使用方式:")
	fmt.Println("  BILIBILI_COOKIES='SESSDATA=...; bili_jct=...' go run examples/comment/main.go add 1 243322853 '测试评论'")
	fmt.Println("  go run examples/comment/main.go like 1 243322853 3039053308")
	fmt.Println("  go run examples/comment/main.go hate 1 243322853 3039053308 1")
	fmt.Println("  go run examples/comment/main.go delete 1 243322853 3039053308")
	fmt.Println("  go run examples/comment/main.go top 1 243322853 2940645593 1")
	fmt.Println("  go run examples/comment/main.go report 1 243322853 3039053308 4")
	fmt.Println("  go run examples/comment/main.go report 1 243322853 3039053308 0 '其他举报原因'")
	fmt.Println("")
	fmt.Println("说明:")
	fmt.Println("  - type=1 表示视频评论区")
	fmt.Println("  - action 默认是 1，传 0 表示取消点赞/取消点踩/取消置顶")
	fmt.Println("  - Cookies 优先读取环境变量 BILIBILI_COOKIES，其次读取当前目录 login.json")
	fmt.Println("  - login.json 可通过登录示例生成")
	fmt.Println("")
	fmt.Println("常见评论区类型:")
	fmt.Println("  - 1: 视频")
	fmt.Println("  - 11: 图文")
	fmt.Println("  - 17: 动态")
	fmt.Println("  - 19: 音频")
	fmt.Println("  - 21: 漫画")
	fmt.Println("  - 22: 活动")
}
