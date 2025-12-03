# 性能测试工具使用说明

## 功能说明

提供三种测试模式来验证并发上传的性能提升。

## 使用方法

### 1. 设置环境变量

```bash
export BILI_SESSDATA="你的SESSDATA"
export BILI_JCT="你的bili_jct"
```

### 2. 运行测试

#### 单线程上传测试

```bash
go run benchmark.go \
  -url "https://your-cos.com/video.mp4" \
  -mode single
```

#### 并发上传测试（默认3并发）

```bash
go run benchmark.go \
  -url "https://your-cos.com/video.mp4" \
  -mode concurrent \
  -c 3
```

#### 完整性能对比测试

```bash
go run benchmark.go \
  -url "https://your-cos.com/video.mp4" \
  -mode compare
```

### 3. 参数说明

| 参数 | 说明 | 必需 | 默认值 |
|------|------|------|--------|
| `-url` | 视频URL地址 | ✅ | - |
| `-name` | 文件名 | ❌ | test-video.mp4 |
| `-mode` | 测试模式 | ❌ | concurrent |
| `-c` | 并发数 | ❌ | 3 |

**测试模式：**
- `single`: 仅测试单线程上传
- `concurrent`: 仅测试并发上传
- `compare`: 完整对比测试（1/2/3/5并发）

## 输出示例

### 单线程测试

```
========================================
🐢 单线程上传测试
========================================
📁 视频信息:
   URL: https://bucket.cos.ap-guangzhou.myqcloud.com/video.mp4
   文件名: test-video.mp4
   大小: 291.11 MB
   模式: single

📤 Uploading chunk 1/30 (3.3%)...
📤 Uploading chunk 2/30 (6.7%)...
...

========================================
✅ 上传成功
========================================
模式: 单线程
文件名: n240728ad1p51if4g3ke
耗时: 3m2s
文件大小: 291.11 MB
平均速度: 1.60 MB/s
========================================
```

### 并发测试

```
========================================
🚀 3并发上传测试
========================================
📁 视频信息:
   URL: https://bucket.cos.ap-guangzhou.myqcloud.com/video.mp4
   文件名: test-video.mp4
   大小: 291.11 MB
   模式: concurrent
   并发数: 3

✅ Selected best endpoint: //upos-cs-upcdntxa.bilivideo.com (latency: 45ms)
✅ [Worker-0] Chunk 1/30 uploaded (3.3%) | Speed: 12.5 MB/s | ETA: 23s
✅ [Worker-1] Chunk 2/30 uploaded (6.7%) | Speed: 13.2 MB/s | ETA: 21s
✅ [Worker-2] Chunk 3/30 uploaded (10.0%) | Speed: 13.8 MB/s | ETA: 19s
...
🎉 All 30 chunks uploaded successfully! Total time: 1m10s, Average speed: 4.16 MB/s

========================================
✅ 上传成功
========================================
模式: 3并发
文件名: n240728ad1p51if4g3ke
耗时: 1m10s
文件大小: 291.11 MB
平均速度: 4.16 MB/s
========================================
```

### 性能对比测试

```
========================================
📊 性能对比测试
========================================

[测试 1/4] 单线程上传...
✅ 单线程完成: 3m2s

[测试 2/4] 2并发上传...
✅ 2并发完成: 1m35s

[测试 3/4] 3并发上传...
✅ 3并发完成: 1m10s

[测试 4/4] 5并发上传...
✅ 5并发完成: 58s

========================================
📈 性能对比结果
========================================

模式       | 耗时         | 速度         | 提升
-----------|--------------|--------------|----------
单线程     | 3m2s         |     1.60 MB/s | 1.00x
2并发      | 1m35s        |     3.07 MB/s | 1.92x
3并发      | 1m10s        |     4.16 MB/s | 2.60x
5并发      | 58s          |     5.02 MB/s | 3.14x

========================================
💡 推荐建议
========================================
🏆 最快模式: 5并发 (耗时 58s)
⚡ 性能提升: 3.14x
```

## 测试建议

### 1. 小文件测试（<100MB）
```bash
# 单线程和3并发差异不明显
go run benchmark.go -url "小文件URL" -mode compare
```

### 2. 中等文件测试（100MB-1GB）
```bash
# 3并发最佳，性能提升明显
go run benchmark.go -url "中等文件URL" -mode concurrent -c 3
```

### 3. 大文件测试（>1GB）
```bash
# 可尝试5并发，但需注意内存占用
go run benchmark.go -url "大文件URL" -mode concurrent -c 5
```

### 4. 完整对比测试
```bash
# 测试所有模式，找出最适合你网络的配置
go run benchmark.go -url "测试文件URL" -mode compare
```

## 注意事项

1. **完整对比测试会多次上传**：`compare` 模式会上传4次，需要足够的时间
2. **网络带宽限制**：如果本地网络带宽有限，增加并发数效果不明显
3. **B站限流**：过高的并发可能触发B站的上传限流
4. **内存占用**：每个并发约占用10MB内存（分块大小）
5. **测试间隔**：完整对比测试会在每次上传后等待2秒

## 故障排查

### 问题：获取文件大小失败
```
❌ 获取文件大小失败: HEAD request failed: 403 Forbidden
```

**解决方案：**
- 检查COS URL是否可公开访问
- 使用带签名的URL
- 确认网络连接正常

### 问题：上传速度慢
```
平均速度: 0.5 MB/s
```

**解决方案：**
1. 检查本地网络带宽
2. 尝试增加并发数
3. 更换COS区域或使用CDN

### 问题：上传失败
```
❌ 上传失败: failed to upload chunks
```

**解决方案：**
1. 检查SESSDATA是否有效
2. 减少并发数重试
3. 查看详细错误日志

## 性能分析

根据测试结果，你可以：

1. **确定最佳并发数**：通过对比测试找出性能最好的配置
2. **评估网络质量**：上传速度反映网络带宽和稳定性
3. **优化上传策略**：根据文件大小选择合适的并发数

## 示例脚本

### 批量测试不同并发数

```bash
#!/bin/bash

VIDEO_URL="https://your-cos.com/test-video.mp4"

for c in 1 2 3 4 5; do
  echo "Testing with concurrency: $c"
  go run benchmark.go -url "$VIDEO_URL" -mode concurrent -c $c
  echo ""
  sleep 5
done
```

### 测试不同大小的文件

```bash
#!/bin/bash

# 小文件
go run benchmark.go -url "https://your-cos.com/small.mp4" -mode compare

# 中等文件
go run benchmark.go -url "https://your-cos.com/medium.mp4" -mode compare

# 大文件
go run benchmark.go -url "https://your-cos.com/large.mp4" -mode compare
```
