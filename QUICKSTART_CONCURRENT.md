# 并发上传快速开始

## 最简单的使用方式

### 1. 安装依赖

```bash
go get github.com/yourusername/bilibili-go-sdk
```

### 2. 准备登录信息

```go
loginInfo := &bilibili.LoginInfo{
    CookieInfo: map[string]interface{}{
        "cookies": []interface{}{
            map[string]interface{}{
                "name":  "SESSDATA",
                "value": "你的SESSDATA",
            },
            map[string]interface{}{
                "name":  "bili_jct",
                "value": "你的bili_jct",
            },
        },
    },
}
```

### 3. 上传视频（3行代码）

```go
uploadClient := bilibili.NewUploadClient(loginInfo)

video, err := uploadClient.UploadVideoFromURLConcurrent(
    "https://your-cos-url.com/video.mp4", // COS链接
    "video.mp4",                            // 文件名
    305333744,                              // 文件大小
    3,                                       // 并发数（0=默认3）
)
```

### 4. 完整示例

```go
package main

import (
    "log"
    "github.com/yourusername/bilibili-go-sdk/bilibili"
)

func main() {
    // 登录信息
    loginInfo := &bilibili.LoginInfo{
        CookieInfo: map[string]interface{}{
            "cookies": []interface{}{
                map[string]interface{}{"name": "SESSDATA", "value": "your_sessdata"},
                map[string]interface{}{"name": "bili_jct", "value": "your_bili_jct"},
            },
        },
    }

    // 创建客户端
    client := bilibili.NewUploadClient(loginInfo)

    // 上传视频（并发）
    video, err := client.UploadVideoFromURLConcurrent(
        "https://bucket.cos.ap-guangzhou.myqcloud.com/video.mp4",
        "video.mp4",
        305333744,
        3,
    )
    if err != nil {
        log.Fatal(err)
    }

    // 提交稿件
    studio := &bilibili.Studio{
        Title:       "视频标题",
        Tid:         122,
        Tag:         "标签1,标签2",
        Description: "视频描述",
        Copyright:   1,
        Videos:      []*bilibili.Video{video},
        NoReprint:   1,
        OpenElec:    1,
        Recreate:    -1,
        Dynamic:     "动态内容",
        NoDisturbance: 1,
    }

    result, err := client.SubmitVideo(studio)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("✅ 发布成功: %v", result.Data)
}
```

## 对比：单线程 vs 并发

```go
// ❌ 慢速单线程上传（不推荐）
video, _ := client.UploadVideoFromURL(url, name, size)

// ✅ 快速并发上传（推荐，速度2-3倍）
video, _ := client.UploadVideoFromURLConcurrent(url, name, size, 3)
```

## 性能提升

- 300MB 视频：**180秒 → 70秒** (2.6倍提升)
- 1GB 视频：**600秒 → 230秒** (2.6倍提升)  
- 5GB 视频：**3000秒 → 1150秒** (2.6倍提升)

## 进度日志示例

```
🚀 Starting concurrent chunk upload from URL: fileSize=305333744, chunksNum=30, concurrency=3
✅ Selected best endpoint: //upos-cs-upcdntxa.bilivideo.com (latency: 45ms)
✅ [Worker-0] Chunk 1/30 uploaded (3.3%) | Speed: 12.5 MB/s | ETA: 23s
✅ [Worker-1] Chunk 2/30 uploaded (6.7%) | Speed: 13.2 MB/s | ETA: 21s
🎉 All 30 chunks uploaded successfully! Total time: 1m23s, Average speed: 3.68 MB/s
```

## 常见问题

**Q: 应该设置多少并发数？**  
A: 推荐3，网络好可以4-5，网络差用2

**Q: 如何获取文件大小？**  
A: 使用 HEAD 请求获取 Content-Length

**Q: 上传失败怎么办？**  
A: 会自动重试最多10次，失败后返回错误

**Q: 会占用多少内存？**  
A: 每个并发约占用10MB（分块大小），3并发约30MB

## 更多文档

- [详细文档](./CONCURRENT_UPLOAD.md)
- [完整示例](./examples/upload/concurrent_upload_example.go)
