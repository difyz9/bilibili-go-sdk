# Bilibili Go SDK - 使用示例

本文档提供了 Bilibili Go SDK 的详细使用示例，涵盖所有主要功能。

## 📋 目录

1. [认证登录](#认证登录)
2. [用户信息](#用户信息)
3. [评论操作](#评论操作)
4. [视频管理](#视频管理)  
5. [标签管理](#标签管理)
6. [数据统计](#数据统计)
7. [直播功能](#直播功能)
8. [创作者工具](#创作者工具)
9. [高级用法](#高级用法)

## 🔐 认证登录

### 二维码登录

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/difyz9/bilibili-go-sdk/bilibili"
)

func qrLogin() {
    client := bilibili.NewClient()
    
    // 获取二维码
    qrResp, err := client.GetQRCode()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("请扫描二维码: %s\n", qrResp.Data.URL)
    
    // 轮询登录状态
    for {
        loginInfo, err := client.PollQRCode(qrResp.Data.AuthCode)
        if err == nil {
            fmt.Printf("登录成功! 用户: %s\n", loginInfo.TokenInfo.Uname)
            cookies := loginInfo.GetCookieString()
            fmt.Printf("Cookies: %s\n", cookies)
            break
        }
        
        fmt.Println("等待扫码...")
        time.Sleep(3 * time.Second)
    }
}
```

### 短信验证码登录

```go
func smsLogin() {
    client := bilibili.NewClient()
    
    // 发送短信验证码
    tel := "13800138000"
    cid := "86" // 中国区号
    
    // 检查手机号是否有效
    valid, err := client.CheckTelValid(tel, cid)
    if err != nil || !valid {
        log.Fatal("手机号无效")
    }
    
    // 发送验证码
    smsResp, err := client.SendSMS(tel, cid)
    if err != nil {
        log.Fatal("发送验证码失败:", err)
    }
    
    // 用户输入验证码
    fmt.Print("请输入验证码: ")
    var code string
    fmt.Scanln(&code)
    
    // 使用验证码登录
    loginReq := &bilibili.SMSLoginRequest{
        Tel:      tel,
        Cid:      cid,
        Code:     code,
        LoginKey: smsResp.CaptchaKey,
    }
    
    loginResp, err := client.SMSLogin(loginReq)
    if err != nil {
        log.Fatal("登录失败:", err)
    }
    
    fmt.Printf("登录成功! Cookies: %s\n", loginResp.CookieInfo.Cookies)
}
```

### Cookie认证

```go
func cookieAuth() {
    // 使用已有的Cookie
    cookiesFile := "cookies.txt"
    cookies := "buvid3=...; bili_jct=...; SESSDATA=..."
    
    // 创建Cookie认证器
    auth := bilibili.NewCookieAuth(cookies)
    
    // 验证Cookie有效性
    if !auth.IsValid() {
        log.Fatal("Cookie无效或已过期")
    }
    
    // 获取用户信息
    userInfo, err := auth.GetUserInfo()
    if err != nil {
        log.Fatal("获取用户信息失败:", err)
    }
    
    fmt.Printf("当前用户: %s (UID: %d)\n", userInfo.Name, userInfo.Mid)
    
    // 保存Cookie到文件
    err = auth.SaveToFile(cookiesFile)
    if err != nil {
        log.Printf("保存Cookie失败: %v", err)
    }
}
```

## � 用户信息

### 获取当前用户信息和粉丝数

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/difyz9/bilibili-go-sdk/bilibili"
)

func getUserInfo(cookies string) {
    client := bilibili.NewClient()
    
    // 获取当前登录用户的详细信息
    myInfo, err := client.GetMyInfoWithRetry(cookies, 3)
    if err != nil {
        log.Fatal("获取用户信息失败:", err)
    }
    
    // 显示用户信息
    fmt.Println("当前登录用户信息:")
    fmt.Printf("  用户ID: %d\n", myInfo.Mid)
    fmt.Printf("  用户名: %s\n", myInfo.Name)
    fmt.Printf("  粉丝数: %d\n", myInfo.Follower)
    fmt.Printf("  关注数: %d\n", myInfo.Following)
    fmt.Printf("  等级: %d\n", myInfo.Level)
    fmt.Printf("  硬币数: %d\n", myInfo.GetCoins())
    fmt.Printf("  签名: %s\n", myInfo.Sign)
    fmt.Printf("  头像: %s\n", myInfo.Face)
}
```

### 获取其他用户的统计信息

```go
func getPublicUserStats(mid int64) {
    client := bilibili.NewClient()
    
    // 获取用户关系统计（粉丝数、关注数等）
    userStat, err := client.GetUserStat(mid)
    if err != nil {
        log.Fatal("获取用户统计失败:", err)
    }
    
    fmt.Printf("用户 %d 的统计信息:\n", mid)
    fmt.Printf("  粉丝数: %d\n", userStat.Follower)
    fmt.Printf("  关注数: %d\n", userStat.Following)
    fmt.Printf("  悄悄关注: %d\n", userStat.Whisper)
    fmt.Printf("  黑名单: %d\n", userStat.Black)
}
```

### 完整示例

运行示例代码：

```bash
# 获取用户信息和粉丝数
go run examples/user_stats/main.go
```

示例代码位置: `examples/user_stats/main.go`

## 评论操作

### 可运行示例

```bash
# 通过环境变量提供 Cookies，发表评论
BILIBILI_COOKIES='SESSDATA=xxx; bili_jct=xxx' go run examples/comment/main.go add 1 243322853 '测试评论'

# 点赞评论
go run examples/comment/main.go like 1 243322853 3039053308

# 举报评论
go run examples/comment/main.go report 1 243322853 3039053308 4
```

示例代码位置: `examples/comment/main.go`

说明：
- `type=1` 表示视频评论区
- 优先读取环境变量 `BILIBILI_COOKIES`，如果未设置则尝试读取当前目录的 `login.json`
- `login.json` 可通过登录示例生成

## �📹 视频管理

### 视频上传

```go
func uploadVideo(cookies string) {
    client := bilibili.NewClient()
    
    // 上传视频文件
    videoPath := "/path/to/video.mp4"
    fmt.Println("开始上传视频...")
    
    uploadResp, err := client.UploadVideo(videoPath, cookies)
    if err != nil {
        log.Fatal("视频上传失败:", err)
    }
    

    status, err := client.WaitForVideoReviewPassed(result.BVID, cookies, 15*time.Second, 30*time.Minute)
    if err != nil {
        log.Printf("审核未通过或等待超时: %v", err)
        if status != nil {
            log.Printf("当前状态: %s", status.StateDesc)
        }
        return
    }

    log.Printf("审核通过: %s", status.BVid)
    fmt.Printf("视频上传成功! 文件名: %s\n", uploadResp.Data.Filename)
    
    // 上传封面
    coverPath := "/path/to/cover.jpg"
    fmt.Println("开始上传封面...")
    
    coverResp, err := client.UploadCover(coverPath, cookies)
    if err != nil {
        log.Fatal("封面上传失败:", err)
    }
    
    fmt.Printf("封面上传成功! URL: %s\n", coverResp.Data.URL)
    
    return uploadResp.Data.Filename, coverResp.Data.URL
}
```

## 🛠️ 创作者工具

### 合集管理

```go
func manageSeason(cookies string) {
    client := bilibili.NewClient()

    seasonID, err := client.CreateSeason(&bilibili.SeasonCreateRequest{
        Title: "系列节目",
        Desc:  "创作中心合集接口示例",
        Cover: "https://i0.hdslb.com/bfs/archive/example.jpg",
    }, cookies)
    if err != nil {
        log.Fatal("创建合集失败:", err)
    }

    seasons, err := client.GetSeasonList(&bilibili.SeasonListParams{Page: 1, PageSize: 30}, cookies)
    if err != nil {
        log.Fatal("获取合集列表失败:", err)
    }

    fmt.Printf("当前共有 %d 个合集\n", seasons.Total)

    err = client.EditSeason(&bilibili.SeasonEditRequest{
        Season: bilibili.SeasonEditInfo{
            ID:    seasonID,
            Title: "系列节目-更新版",
            Cover: "https://i0.hdslb.com/bfs/archive/example.jpg",
            Desc:  "更新后的简介",
        },
    }, cookies)
    if err != nil {
        log.Fatal("编辑合集失败:", err)
    }

    if err := client.DeleteSeason(seasonID, cookies); err != nil {
        log.Fatal("删除合集失败:", err)
    }
}
```

### 视频发布

```go
func publishVideo(filename, coverURL, cookies string) {
    client := bilibili.NewClient()
    
    // 构建发布请求
    submitReq := &bilibili.SubmitRequest{
        Cover:     coverURL,
        Title:     "我的第一个视频",
        Desc:      "这是一个使用Go SDK上传的测试视频",
        Tag:       "编程,Go语言,教程",
        Tid:       174, // 计算机技术分区
        Copyright: 1,   // 自制
        Source:    "",  // 自制视频无需来源
        Dynamic:   "发布了新视频，欢迎观看~",
        Videos: []bilibili.VideoMeta{
            {
                Filename: filename,
                Title:    "P1",
                Desc:     "",
            },
        },
        OpenElec:  1, // 开启充电
        NoReprint: 0, // 允许转载
    }
    
    fmt.Println("开始发布视频...")
    submitResp, err := client.SubmitVideo(submitReq, cookies)
    if err != nil {
        log.Fatal("视频发布失败:", err)
    }
    
    fmt.Printf("视频发布成功!\n")
    fmt.Printf("AV号: %d\n", submitResp.Data.AID)
    fmt.Printf("BV号: %s\n", submitResp.Data.BVid)
}
```

### 视频编辑

```go
func editVideo(aid int64, cookies string) {
    client := bilibili.NewClient()
    
    // 先获取视频信息
    bvid := "BV1xx411c7mD" // 示例BV号
    videoInfo, err := client.GetVideoInfo(bvid, cookies)
    if err != nil {
        log.Fatal("获取视频信息失败:", err)
    }
    
    fmt.Printf("当前标题: %s\n", videoInfo.Title)
    
    // 修改视频信息
    editReq := &bilibili.VideoEditRequest{
        AID:         aid,
        BVid:        bvid,
        Title:       "新的视频标题",
        Desc:        "更新后的视频描述",
        Cover:       videoInfo.Cover,
        Tag:         "新标签1,新标签2,新标签3",
        Tid:         videoInfo.Tid,
        Copyright:   videoInfo.Copyright,
        Source:      videoInfo.Source,
        OpenElec:    1,
        NoReprint:   0,
        SelectiOn:   0,
        Dynamic:     "更新了视频信息",
        Interactive: 0,
    }
    
    err = client.EditVideo(editReq, cookies)
    if err != nil {
        log.Fatal("编辑视频失败:", err)
    }
    
    fmt.Println("视频信息更新成功!")
}
```

### 获取视频列表

```go
func getMyVideos(cookies string) {
    client := bilibili.NewClient()
    
    page := 1
    pageSize := 10
    
    videos, err := client.GetMyVideos(page, pageSize, cookies)
    if err != nil {
        log.Fatal("获取视频列表失败:", err)
    }
    
    fmt.Printf("共找到 %d 个视频:\n", len(videos))
    for i, video := range videos {
        fmt.Printf("%d. %s (AV%d, %s)\n", 
            i+1, video.Title, video.AID, video.BVid)
        fmt.Printf("   状态: %d, 播放: %d, 时长: %d秒\n", 
            video.State, video.Duration, video.Duration)
    }
}
```

### 查询投稿审核状态

```go
func getVideoReviewStatus(cookies string) {
    client := bilibili.NewClient()

    status, err := client.GetVideoReviewStatus("BV1xx411c7mD", cookies)
    if err != nil {
        log.Fatal("获取投稿审核状态失败:", err)
    }

    fmt.Printf("是否通过审核: %v\n", status.Passed)
    fmt.Printf("当前状态: %s\n", status.StateDesc)
    if status.RejectReason != "" {
        fmt.Printf("驳回原因: %s\n", status.RejectReason)
    }
}
```

## 🏷️ 标签管理

### 标签验证和推荐

```go
func tagManagement(cookies string) {
    client := bilibili.NewClient()
    
    // 检查单个标签
    tag := "编程"
    valid, err := client.CheckTag(tag)
    if err != nil {
        log.Printf("检查标签失败: %v", err)
    } else {
        fmt.Printf("标签'%s'有效性: %v\n", tag, valid)
    }
    
    // 获取推荐标签
    title := "Go语言入门教程"
    desc := "从零开始学习Go语言编程"
    
    tags, err := client.GetRecommendedTags(title, desc, cookies)
    if err != nil {
        log.Printf("获取推荐标签失败: %v", err)
    } else {
        fmt.Printf("推荐标签 (%d个):\n", len(tags))
        for _, tag := range tags {
            fmt.Printf("- %s (类型: %d)\n", tag.Name, tag.Type)
        }
    }
    
    // 批量验证标签
    tagList := []string{"编程", "Go语言", "教程", "技术", "无效标签123"}
    results, err := client.ValidateTags(tagList)
    if err != nil {
        log.Printf("批量验证失败: %v", err)
    } else {
        fmt.Println("批量验证结果:")
        for tag, valid := range results {
            status := "有效"
            if !valid {
                status = "无效"
            }
            fmt.Printf("- %s: %s\n", tag, status)
        }
    }
    
    // 格式化标签
    rawTags := []string{"编程", "Go语言", "教程", "技术", "开源", "后端", "微服务", "云原生", "容器", "k8s", "超长标签会被过滤"}
    formatted := bilibili.FormatTags(rawTags)
    fmt.Printf("格式化后的标签: %s\n", formatted)
}
```

## 📊 数据统计

### 视频统计信息

```go
func getVideoStats() {
    client := bilibili.NewClient()
    
    bvid := "BV1xx411c7mD"
    
    // 获取视频统计信息
    stat, err := client.GetVideoStat(bvid)
    if err != nil {
        log.Fatal("获取视频统计失败:", err)
    }
    
    fmt.Printf("视频统计信息 (BV: %s):\n", stat.BVid)
    fmt.Printf("播放量: %d\n", stat.View)
    fmt.Printf("弹幕数: %d\n", stat.Danmaku)
    fmt.Printf("评论数: %d\n", stat.Reply)
    fmt.Printf("点赞数: %d\n", stat.Like)
    fmt.Printf("投币数: %d\n", stat.Coin)
    fmt.Printf("收藏数: %d\n", stat.Favorite)
    fmt.Printf("分享数: %d\n", stat.Share)
}
```

### UP主统计信息

```go
func getUpStats(cookies string) {
    client := bilibili.NewClient()
    
    // 获取我的UP主统计
    upStat, err := client.GetMyUpStat(cookies)
    if err != nil {
        log.Fatal("获取UP主统计失败:", err)
    }
    
    fmt.Println("我的UP主统计:")
    fmt.Printf("视频总播放量: %d\n", upStat.Archive.View)
    fmt.Printf("专栏总阅读量: %d\n", upStat.Article.View)
    fmt.Printf("获赞总数: %d\n", upStat.Likes)
    
    // 获取其他用户统计
    mid := int64(12345678) // 目标用户UID
    userStat, err := client.GetUserStat(mid)
    if err != nil {
        log.Printf("获取用户统计失败: %v", err)
    } else {
        fmt.Printf("\n用户 %d 的统计:\n", mid)
        fmt.Printf("关注数: %d\n", userStat.Following)
        fmt.Printf("粉丝数: %d\n", userStat.Follower)
    }
}
```

### 视频详细分析

```go
func getVideoAnalytics(aid int64, cookies string) {
    client := bilibili.NewClient()
    
    // 获取视频详细分析数据 (需要是视频作者)
    analytics, err := client.GetVideoAnalytics(aid, cookies)
    if err != nil {
        log.Fatal("获取视频分析失败:", err)
    }
    
    fmt.Println("视频详细分析数据:")
    for key, value := range analytics {
        fmt.Printf("%s: %v\n", key, value)
    }
    
    // 获取趋势数据
    startDate := time.Now().AddDate(0, 0, -7) // 7天前
    endDate := time.Now()
    
    trendData, err := client.GetTrendData(aid, startDate, endDate, cookies)
    if err != nil {
        log.Printf("获取趋势数据失败: %v", err)
    } else {
        fmt.Println("7天趋势数据:")
        for key, value := range trendData {
            fmt.Printf("%s: %v\n", key, value)
        }
    }
}
```

## 🎥 直播功能

### 直播管理

```go
func liveManagement(roomID int64, cookies string) {
    client := bilibili.NewClient()
    
    // 获取直播间信息
    roomInfo, err := client.GetLiveRoomInfo(roomID)
    if err != nil {
        log.Fatal("获取直播间信息失败:", err)
    }
    
    fmt.Printf("直播间 %d 信息:\n", roomID)
    fmt.Printf("短号: %d\n", roomInfo.ShortID)
    fmt.Printf("主播UID: %d\n", roomInfo.UID)
    fmt.Printf("直播状态: %d\n", roomInfo.LiveStatus)
    fmt.Printf("是否隐藏: %v\n", roomInfo.IsHidden)
    
    // 如果是自己的直播间，可以进行管理操作
    if roomInfo.UID == getCurrentUserUID() {
        // 更新直播标题
        err = client.UpdateLiveTitle(roomID, "新的直播标题", cookies)
        if err != nil {
            log.Printf("更新直播标题失败: %v", err)
        } else {
            fmt.Println("直播标题更新成功")
        }
        
        // 开始直播
        area := 371 // 聊天分区
        err = client.StartLive(roomID, area, cookies)
        if err != nil {
            log.Printf("开始直播失败: %v", err)
        } else {
            fmt.Println("直播开始成功")
        }
        
        // 等待一段时间后停止直播
        time.Sleep(10 * time.Second)
        
        err = client.StopLive(roomID, cookies)
        if err != nil {
            log.Printf("停止直播失败: %v", err)
        } else {
            fmt.Println("直播已停止")
        }
    }
}
```

### 获取直播流信息

```go
func getLiveStream(roomID int64, cookies string) {
    client := bilibili.NewClient()
    
    quality := 10000 // 原画质量
    
    streamInfo, err := client.GetLiveStreamInfo(roomID, quality, cookies)
    if err != nil {
        log.Fatal("获取直播流信息失败:", err)
    }
    
    fmt.Printf("直播流信息 (房间: %d):\n", roomID)
    fmt.Printf("当前画质: %d (%s)\n", 
        streamInfo.CurrentQuality, streamInfo.CurrentQualityName)
    
    fmt.Printf("可用画质: %v\n", streamInfo.AcceptQuality)
    
    if len(streamInfo.Durl) > 0 {
        fmt.Printf("流地址: %s\n", streamInfo.Durl[0].URL)
    }
}
```

## 🛠️ 创作者工具

### 草稿管理

```go
func draftManagement(cookies string) {
    client := bilibili.NewClient()
    
    // 获取草稿列表
    page := 1
    pageSize := 10
    
    drafts, err := client.GetDraftList(page, pageSize, cookies)
    if err != nil {
        log.Fatal("获取草稿列表失败:", err)
    }
    
    fmt.Printf("草稿列表 (%d个):\n", len(drafts))
    for i, draft := range drafts {
        fmt.Printf("%d. %s (ID: %d)\n", i+1, draft.Title, draft.ID)
        fmt.Printf("   状态: %d, 创建时间: %d\n", draft.Status, draft.Created)
    }
    
    // 创建新草稿
    draftID, err := client.SaveDraft(
        "草稿标题",
        "草稿描述",
        "https://example.com/cover.jpg",
        "标签1,标签2",
        174, // 分区ID
        cookies,
    )
    if err != nil {
        log.Printf("保存草稿失败: %v", err)
    } else {
        fmt.Printf("草稿保存成功，ID: %d\n", draftID)
        
        // 删除草稿
        err = client.DeleteDraft(draftID, cookies)
        if err != nil {
            log.Printf("删除草稿失败: %v", err)
        } else {
            fmt.Println("草稿删除成功")
        }
    }
}
```

### 模板管理

```go
func templateManagement(cookies string) {
    client := bilibili.NewClient()
    
    // 获取模板列表
    templates, err := client.GetTemplateList(1, 10, cookies)
    if err != nil {
        log.Fatal("获取模板列表失败:", err)
    }
    
    fmt.Printf("模板列表 (%d个):\n", len(templates))
    for i, template := range templates {
        fmt.Printf("%d. %s (ID: %d)\n", i+1, template.Name, template.ID)
        fmt.Printf("   描述: %s\n", template.Description)
        fmt.Printf("   标签: %s\n", template.Tag)
    }
    
    // 创建新模板
    templateID, err := client.CreateTemplate(
        "我的模板",
        "这是一个模板描述",
        "https://example.com/template_cover.jpg",
        "模板,测试",
        174, // 分区ID
        cookies,
    )
    if err != nil {
        log.Printf("创建模板失败: %v", err)
    } else {
        fmt.Printf("模板创建成功，ID: %d\n", templateID)
    }
}
```

## 🚀 高级用法

### 错误处理

```go
func advancedErrorHandling() {
    client := bilibili.NewClient()
    
    // 使用重试机制
    var loginInfo *bilibili.LoginInfo
    var err error
    
    for retries := 0; retries < 3; retries++ {
        qrResp, err := client.GetQRCode()
        if err != nil {
            log.Printf("获取二维码失败 (重试 %d/3): %v", retries+1, err)
            time.Sleep(time.Duration(retries+1) * time.Second)
            continue
        }
        
        fmt.Printf("请扫描二维码: %s\n", qrResp.Data.URL)
        
        for polls := 0; polls < 60; polls++ { // 最多轮询3分钟
            loginInfo, err = client.PollQRCode(qrResp.Data.AuthCode)
            if err == nil {
                fmt.Println("登录成功!")
                return
            }
            
            // 根据错误类型决定是否继续
            if strings.Contains(err.Error(), "expired") {
                fmt.Println("二维码已过期，重新获取...")
                break
            }
            
            time.Sleep(3 * time.Second)
        }
    }
    
    log.Fatal("登录失败，已超过最大重试次数")
}
```

### 并发处理

```go
func concurrentProcessing(videoList []string, cookies string) {
    client := bilibili.NewClient()
    
    // 使用goroutine并发处理多个视频
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 3) // 限制并发数为3
    
    for _, bvid := range videoList {
        wg.Add(1)
        go func(vid string) {
            defer wg.Done()
            
            semaphore <- struct{}{} // 获取信号量
            defer func() { <-semaphore }() // 释放信号量
            
            stat, err := client.GetVideoStat(vid)
            if err != nil {
                log.Printf("获取视频 %s 统计失败: %v", vid, err)
                return
            }
            
            fmt.Printf("视频 %s: 播放 %d, 点赞 %d\n", 
                vid, stat.View, stat.Like)
        }(bvid)
    }
    
    wg.Wait()
    fmt.Println("所有视频处理完成")
}
```

### 自定义HTTP客户端

```go
func customHTTPClient() {
    // 创建自定义HTTP客户端
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     30 * time.Second,
        },
    }
    
    // 使用自定义配置创建客户端
    config := &bilibili.Config{
        UserAgent:  "自定义UserAgent",
        HTTPClient: httpClient,
    }
    
    client := bilibili.NewClientWithConfig(config)
    
    // 使用客户端进行操作
    qrResp, err := client.GetQRCode()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("二维码URL: %s\n", qrResp.Data.URL)
}
```

## 🔧 配置和最佳实践

### Cookie管理

```go
func cookieManagement() {
    // 从文件加载Cookie
    cookieFile := "bilibili_cookies.txt"
    auth, err := bilibili.LoadCookieFromFile(cookieFile)
    if err != nil {
        log.Printf("加载Cookie失败: %v", err)
        // 重新登录获取Cookie
        // ...
        return
    }
    
    // 验证Cookie是否仍然有效
    if !auth.IsValid() {
        log.Println("Cookie已过期，需要重新登录")
        return
    }
    
    // 使用Cookie获取用户信息
    userInfo, err := auth.GetUserInfo()
    if err != nil {
        log.Printf("获取用户信息失败: %v", err)
        return
    }
    
    fmt.Printf("当前用户: %s\n", userInfo.Name)
    
    // 定期刷新Cookie (如果支持)
    refreshed := auth.Refresh()
    if refreshed {
        // 保存更新后的Cookie
        auth.SaveToFile(cookieFile)
        fmt.Println("Cookie已刷新")
    }
}
```

### 日志记录

```go
func setupLogging() {
    // 配置日志输出
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    // 创建日志文件
    logFile, err := os.OpenFile("bilibili_sdk.log", 
        os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("创建日志文件失败:", err)
    }
    
    // 同时输出到控制台和文件
    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)
    
    log.Println("日志系统初始化完成")
}
```

## 📝 注意事项

1. **频率限制**: 请遵守B站的API调用频率限制，避免请求过于频繁
2. **Cookie安全**: 妥善保管登录Cookie，不要泄露给他人
3. **WBI签名**: SDK自动处理WBI签名，但需要网络连接获取最新密钥
4. **错误处理**: 建议对所有API调用进行适当的错误处理
5. **版本兼容**: 定期更新SDK以获得最新功能和修复

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个SDK！

## 📄 许可证

MIT License