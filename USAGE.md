# Bilibili Go SDK 使用指南

## 项目结构

```
bilibili-go-sdk/
├── bilibili/           # 核心SDK包
│   ├── client.go       # 客户端核心
│   ├── config.go       # 配置管理
│   ├── types.go        # 数据类型定义
│   ├── auth.go         # 认证相关API
│   ├── upload.go       # 上传核心功能
│   ├── upload_helpers.go # 上传辅助函数
│   ├── upload_types.go # 上传相关类型
│   ├── submit.go       # 投稿和封面上传
│   └── client_test.go  # 单元测试
├── examples/           # 使用示例
│   ├── login/          # 登录示例
│   ├── subtitle/       # 字幕上传与语言映射验证示例
│   ├── upload/         # 上传示例
│   └── complete/       # 完整流程示例
├── go.mod              # Go模块定义
├── README.md           # 项目说明
├── LICENSE             # 许可证
└── USAGE.md            # 本文件
```

## 快速开始

### 1. 安装SDK

```bash
go mod init your-project
go get github.com/difyz9/bilibili-go-sdk
```

### 2. 基本使用

```go
import "github.com/difyz9/bilibili-go-sdk/bilibili"

// 创建客户端
client := bilibili.NewClient()

// 登录流程
qrResp, err := client.GetQRCode()
// 显示二维码让用户扫描
loginInfo, err := client.PollQRCode(qrResp.Data.AuthCode)

// 创建上传客户端
client := bilibili.NewClient()
uploader := bilibili.NewUploadClient(loginInfo)
cookies := loginInfo.GetCookieString()

// 上传视频
video, err := uploader.UploadVideo("/path/to/video.mp4")

// 预测分区
predictions, err := uploader.PredictArchiveTypes(&bilibili.ArchiveTypePredictRequest{
    Filename: video.Filename,
    Title:    "视频标题",
})

tid := 174
if err == nil && len(predictions) > 0 {
    tid = predictions[0].ID
}

// 推荐标签
tags, err := client.RecommendTags(&bilibili.TagRecommendRequest{
    SubtypeID:   tid,
    Title:       "视频标题",
    Filename:    video.Filename,
    Description: "视频描述",
}, cookies)

tagNames := []string{"标签1", "标签2"}
if err == nil && len(tags) > 0 {
    tagNames = []string{tags[0].Name}
}

// 投稿
studio := &bilibili.Studio{
    Title:        "视频标题",
    Desc:         "视频描述",
    Tid:          tid,
    Tag:          bilibili.FormatTags(tagNames),
    DescFormatId: 9999,
    Recreate:     -1,
    WebOS:        3,
    Videos:       []bilibili.Video{*video},
}
result, err := uploader.SubmitVideo(studio)
```

### 3. 运行示例

```bash
# 登录示例
cd examples/login && go run main.go

# 完整流程示例
cd examples/complete
go run main.go login              # 先登录
go run main.go upload video.mp4   # 上传视频

# 字幕上传
go run examples/subtitle/main.go inspect
go run examples/subtitle/main.go upload bili_bvid ./examples/subtitle/FM5-R4VPArw.zh.srt zh ./examples/login_info.json
```

## 核心功能

### 认证模块 (auth.go)
- `GetQRCode()` - 获取登录二维码
- `PollQRCode()` - 轮询登录状态
- `GetMyInfo()` - 获取用户详细信息
- `GetUserBasicInfo()` - 获取指定用户空间详细信息
- `GetArchivePre()` - 获取投稿分区信息

### 评论模块 (comment.go)
- `AddComment()` - 发表评论
- `LikeComment()` - 点赞或取消点赞评论
- `HateComment()` - 点踩或取消点踩评论
- `DeleteComment()` - 删除评论
- `TopComment()` - 置顶或取消置顶评论
- `ReportComment()` - 举报评论

### 上传模块 (upload.go)
- `UploadVideo()` - 从本地文件上传视频
- `UploadVideoFromURL()` - 从URL上传视频
- `UploadCover()` - 上传封面图片
- `SubmitVideo()` - 使用 Web 接口提交视频投稿

### 创作中心投稿辅助接口
- `PredictArchiveTypes()` - 预测稿件分区
- `GetHumanTypeList()` - 获取新分区列表
- `GetUploadTemplates()` - 获取上传模板列表
- `UpdateUploadTemplate()` - 编辑上传模板
- `QueryTopics()` - 查询话题
- `SearchTopics()` - 搜索话题
- `GetArchiveDescFormat()` - 获取简介格式信息
- `ProbeUploadLines()` - 获取上传线路
- `RecommendTags()` - 获取推荐标签
- `CheckTag()` - 校验标签是否可用

### 字幕模块 (subtitle.go)
- `UploadSubtitle()` - 读取 SRT、转换为 BCC 并提交字幕草稿
- `LoadSRTAsBCC()` - 将 SRT 文件转换为 BCC 字幕结构
- `NormalizeSubtitleLanguage()` - 将常见语言别名标准化为接口接受的 `lan`

### 配置模块 (config.go)
- `WithTimeout()` - 设置超时时间
- `WithUserAgent()` - 设置User-Agent
- `WithHTTPClient()` - 自定义HTTP客户端
- `WithProxy()` - 设置代理

## 高级用法

### 自定义配置

```go
client := bilibili.NewClient(
    bilibili.WithTimeout(60 * time.Second),
    bilibili.WithUserAgent("MyApp/1.0"),
)
```

### 错误处理

```go
if bilibili.IsRateLimitError(err) {
    // 处理限流错误
    time.Sleep(time.Minute)
    // 重试
}

if bilibili.IsNetworkError(err) {
    // 处理网络错误
    // SDK内置重试机制会自动处理
}
```

### 断点续传

SDK内置支持分块上传，自动处理网络中断和重试。

### 字幕上传

字幕上传走的是 Bilibili 当前可用的 `draft/save` 流程，SDK 会自动把 SRT 转成接口需要的 BCC JSON。

推荐步骤：

```bash
# 1. 登录并保存会话
go run ./examples/login/main.go

# 2. 检查语言映射
go run ./examples/subtitle/main.go inspect zh zh-Hans zh-TW en-US

# 3. 上传字幕
go run ./examples/subtitle/main.go upload BV16PoVBjE21 ./examples/subtitle/FM5-R4VPArw.zh.srt zh ./examples/login_info.json
```

SDK 调用示例：

```go
client := bilibili.NewClient()
uploader := bilibili.NewSubtitleUploader(client, loginInfo)

if err := uploader.UploadSubtitle("BV16PoVBjE21", "./examples/subtitle/FM5-R4VPArw.zh.srt", "zh"); err != nil {
    log.Fatal(err)
}
```

说明：

- 上传命令参数为：`upload <bvid> <subtitle.srt> <language> [login_info.json]`
- 简体中文字幕建议使用 `zh`；`zh-CN`、`zh-Hans`、`cmn-Hans` 会自动归一化到 `zh`
- 繁体中文字幕会归一化到 `zh-TW`
- 若未指定登录信息文件，示例默认读取 `./examples/login_info.json`

## 注意事项

1. **登录信息保存**: 登录后的 `LoginInfo` 建议保存到文件，避免重复登录
2. **限流处理**: B站API有限流，SDK内置了重试机制
3. **文件大小**: 建议单个视频文件不超过4GB
4. **分区选择**: 投稿时需选择正确的分区ID

## 分区ID参考

常用分区ID：
- 1: 动画
- 3: 音乐  
- 4: 游戏
- 36: 科技
- 160: 生活
- 174: 生活区
- 188: 科技区

获取完整分区列表：
```go
// 内置完整视频分区字典，无需 Cookie
zones := bilibili.GetVideoZones()

// 按 tid 查询分区；若同一个 tid 同时存在重定向条目，会优先返回有效分区
zone, ok := bilibili.GetVideoZoneByTID(176)

// 如果需要保留重定向/重复 tid 的所有匹配项，可使用：
allMatches := bilibili.FindVideoZonesByTID(176)

// 投稿可用分区仍然建议结合 archive/pre 接口获取
archiveData, err := client.GetArchivePre(cookies)
```

## 错误排查

### 登录失败
- 检查网络连接
- 确认二维码未过期
- 验证手机端bilibili应用版本

### 上传失败
- 检查文件是否存在且可读
- 验证登录状态是否有效
- 检查网络连接稳定性

### 投稿失败
- 确认必填字段完整
- 检查分区ID是否正确
- 验证视频格式和大小

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License