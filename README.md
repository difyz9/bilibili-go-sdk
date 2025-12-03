# Bilibili API Go SDK

Bilibili API Go SDK，提供哔哩哔哩平台的全方位功能封装。

## 🌟 功能特性

### 🔐 认证登录
- **二维码登录** - QR码扫码登录，安全便捷
- **短信登录** - 手机号验证码登录  
- **密码登录** - 用户名密码登录
- **Cookie认证** - 基于已有Cookie的认证
- **WBI签名** - 自动处理B站WBI签名验证

### 📊 视频管理
- **视频上传** - 支持分片上传和断点续传
- **🚀 并发上传（NEW）** - 3线程并发上传，速度提升2-3倍
- **从URL上传** - 直接从腾讯COS等云存储URL上传
- **智能Endpoint选择** - 自动选择最优上传节点
- **封面上传** - 视频封面图片上传
- **视频投稿** - 完整的视频发布流程
- **视频编辑** - 修改视频信息（标题、简介、标签等）
- **视频删除** - 删除已发布视频
- **视频列表** - 获取个人视频列表

### 🏷️ 标签管理
- **标签验证** - 检查标签是否有效
- **推荐标签** - 基于内容获取推荐标签
- **标签搜索** - 搜索相关标签
- **批量验证** - 批量验证多个标签
- **标签格式化** - 自动格式化为B站要求格式

### 📊 数据统计
- **视频统计** - 播放量、点赞、收藏等数据
- **用户统计** - 粉丝、关注、获赞等统计
- **UP主统计** - 创作者数据概览
- **详细分析** - 视频详细分析数据（需登录）
- **趋势数据** - 时间段内的数据趋势

### 🎥 直播功能
- **直播间信息** - 获取直播间详细信息
- **直播流信息** - 获取直播流地址和画质
- **开始直播** - 启动直播推流
- **停止直播** - 结束直播
- **更新标题** - 修改直播间标题

### 🛠️ 创作者工具
- **创作者信息** - 获取创作者基本信息
- **草稿管理** - 保存、删除、列表草稿
- **模板管理** - 创建和使用视频模板
- **批量操作** - 批量处理视频内容

### � 用户信息
- **基本信息** - 获取用户基本资料
- **详细信息** - 获取用户详细信息
- **个人空间** - 访问用户空间数据

## 🚀 快速开始

### 安装

```bash
go get github.com/difyz9/bilibili-go-sdk
```

### 基本使用

#### 1. 二维码登录

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/difyz9/bilibili-go-sdk/bilibili"
)

func main() {
    // 创建客户端
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
        
        time.Sleep(3 * time.Second)
    }
}
```

#### 2. 短信登录

```go
// 发送短信验证码
smsResp, err := client.SendSMS("13800138000", "86")
if err != nil {
    log.Fatal(err)
}

// 使用验证码登录
loginReq := &bilibili.SMSLoginRequest{
    Tel:      "13800138000",
    Cid:      "86",
    Code:     "123456", // 用户输入的验证码
    LoginKey: smsResp.CaptchaKey,
}

loginResp, err := client.SMSLogin(loginReq)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("登录成功! Cookies: %s\n", loginResp.CookieInfo.Cookies)
```

#### 3. Cookie认证

```go
// 使用已有的Cookie进行认证
cookies := "buvid3=...; bili_jct=...; SESSDATA=..."
auth := bilibili.NewCookieAuth(cookies)

if auth.IsValid() {
    fmt.Println("Cookie认证成功")
    userInfo, _ := auth.GetUserInfo()
    fmt.Printf("当前用户: %s\n", userInfo.Name)
}
```
    
    fmt.Printf("登录成功! 用户: %s\n", loginInfo.TokenInfo.Uname)
    
    // 获取用户信息
    userInfo, err := client.GetMyInfo(loginInfo.GetCookieString())
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("用户详情: %s (等级 %d)\n", userInfo.Uname, userInfo.Level)
}
```

### 视频上传示例

```go
package main

import (
    "log"
    
    "github.com/difyz9/bilibili-go-sdk/bilibili"
)

func main() {
    // 首先需要登录获取 loginInfo
    // ... (登录代码省略)
    
    // 创建上传客户端
    uploader := bilibili.NewUploadClient(loginInfo)
    
    // 上传视频文件
    video, err := uploader.UploadVideo("/path/to/your/video.mp4")
    if err != nil {
        log.Fatal(err)
    }
    
    // 上传封面
    coverURL, err := uploader.UploadCover("/path/to/cover.jpg")
    if err != nil {
        log.Fatal(err)
    }
    
    // 构建投稿信息
    studio := &bilibili.Studio{
        Title:     "我的视频标题",
        Desc:      "视频描述内容",
        Tid:       174, // 分区ID（生活区）
        Cover:     coverURL,
        Tag:       "标签1,标签2,标签3",
        Copyright: 1, // 原创
        Videos:    []bilibili.Video{*video},
    }
    
    // 提交投稿
    result, err := uploader.SubmitVideo(studio)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("投稿成功! 结果: %+v", result)
}
```

## API 文档

### 客户端 (Client)

#### 创建客户端
```go
client := bilibili.NewClient()
```

#### QR 码登录
```go
// 获取 QR 码
qrResp, err := client.GetQRCode()

// 轮询登录状态
loginInfo, err := client.PollQRCode(authCode)
```

#### 用户信息
```go
// 获取详细用户信息 (推荐)
myInfo, err := client.GetMyInfo(cookies)

// 带重试机制的获取用户信息
myInfo, err := client.GetMyInfoWithRetry(cookies, 3)
```

#### 分区信息
```go
// 获取投稿分区列表
archiveData, err := client.GetArchivePre(cookies)
```

### 上传客户端 (UploadClient)

#### 创建上传客户端
```go
uploader := bilibili.NewUploadClient(loginInfo)
```

#### 🚀 并发视频上传（推荐）

```go
// 从URL并发上传视频（速度提升2-3倍）
video, err := uploader.UploadVideoFromURLConcurrent(
    "https://your-cos.com/video.mp4",  // 视频URL
    "video.mp4",                         // 文件名
    305333744,                           // 文件大小（字节）
    3,                                    // 并发数（推荐3，0=默认）
)

// 特点：
// - 3个goroutine并发上传不同分块
// - 自动选择最优上传节点
// - 实时显示速度和ETA
// - 300MB视频：180秒 → 70秒（2.6倍提升）
```

**详细文档：**
- [并发上传详细指南](./CONCURRENT_UPLOAD.md)
- [快速开始](./QUICKSTART_CONCURRENT.md)
- [完整示例](./examples/upload/concurrent_upload_example.go)

#### 视频上传
```go
// 从本地文件上传
video, err := uploader.UploadVideo("/path/to/video.mp4")

// 从 URL 上传
video, err := uploader.UploadVideoFromURL("https://example.com/video.mp4", "video.mp4", fileSize)
```

#### 封面上传
```go
coverURL, err := uploader.UploadCover("/path/to/cover.jpg")
```

#### 视频投稿
```go
studio := &bilibili.Studio{
    Title:     "视频标题",
    Desc:      "视频描述", 
    Tid:       174,
    Cover:     coverURL,
    Tag:       "标签1,标签2",
    Copyright: 1,
    Videos:    []bilibili.Video{*video},
}

result, err := uploader.SubmitVideo(studio)
```

## 数据结构

### LoginInfo - 登录信息
```go
type LoginInfo struct {
    CookieInfo map[string]interface{} `json:"cookie_info"`
    SSO        []string               `json:"sso"`
    TokenInfo  TokenInfo              `json:"token_info"`
    Platform   string                 `json:"platform"`
}
```

### UserInfo - 用户信息
```go
type MyInfoResponse struct {
    Mid       int64  `json:"mid"`
    Uname     string `json:"uname"`
    Sign      string `json:"sign"`
    Face      string `json:"face"`
    Level     int    `json:"level"`
    Coins     int    `json:"coins"`
    Fans      int    `json:"fans"`
    // ... 更多字段
}
```

### Studio - 投稿信息
```go
type Studio struct {
    Title     string  `json:"title"`
    Desc      string  `json:"desc"`
    Tid       int     `json:"tid"`
    Cover     string  `json:"cover"`
    Tag       string  `json:"tag"`
    Copyright int     `json:"copyright"`
    Videos    []Video `json:"videos"`
    // ... 更多字段
}
```

## 配置选项

### 自定义 HTTP 客户端
```go
import "net/http"

client := bilibili.NewClient()
// 访问内部 httpClient 进行自定义配置
client.SetTimeout(60 * time.Second)
```

### 设置代理
```go
// 通过环境变量
os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")

// 或通过自定义 Transport
transport := &http.Transport{
    Proxy: http.ProxyURL(proxyURL),
}
```

## 错误处理

SDK 提供了详细的错误信息和重试机制：

```go
// 检查是否是限流错误
if bilibili.IsRateLimitError(err) {
    // 处理限流情况
    time.Sleep(time.Minute)
    // 重试
}

// 检查是否是网络错误
if bilibili.IsNetworkError(err) {
    // 处理网络错误
}
```

## 📚 示例代码

SDK 提供了多个完整的示例代码，位于 `examples/` 目录下：

- **`examples/login/`** - 登录示例（二维码登录）
- **`examples/user_stats/`** - 获取用户信息和粉丝数
- **`examples/upload/`** - 视频上传示例
- **`examples/complete/`** - 完整流程示例（登录、上传、投稿）

运行示例：

```bash
# 登录示例
go run examples/login/main.go

# 获取用户信息和粉丝数
go run examples/user_stats/main.go

# 视频上传示例
go run examples/upload/main.go

# 完整流程示例
go run examples/complete/main.go login
go run examples/complete/main.go upload video.mp4 cover.jpg
```

更多示例和详细文档，请查看 [EXAMPLES.md](EXAMPLES.md)

## 注意事项

1. **Cookie 管理**: 登录后的 Cookie 需要妥善保存，用于后续 API 调用
2. **限流处理**: B站 API 有限流机制，SDK 内置了重试逻辑
3. **文件大小**: 视频文件建议不超过 4GB
4. **分区选择**: 投稿时需要选择正确的分区 ID
5. **安全性**: 请勿在公开代码中硬编码敏感信息

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v1.0.0
- 初始发布
- 支持 QR 码登录
- 支持用户信息获取  
- 支持视频上传和投稿
- 支持封面上传
- 内置重试机制