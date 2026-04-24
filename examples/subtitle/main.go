package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

const defaultLoginInfoPath = "login_info.json"

var defaultLoginInfoCandidates = []string{
	defaultLoginInfoPath,
	filepath.Join("examples", "login", defaultLoginInfoPath),
	filepath.Join("examples", "complete", defaultLoginInfoPath),
}

var subtitleLanguageSamples = []string{
	"zh",
	"en",
	"ja",
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "inspect":
		runInspect(os.Args[2:])
	case "upload":
		runUpload(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("字幕上传字段验证示例")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run examples/subtitle/main.go inspect [语言值...]")
	fmt.Println("  go run examples/subtitle/main.go upload <bvid> <subtitle.srt> <language> [login_info.json]")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run examples/subtitle/main.go inspect")
	fmt.Println("  go run examples/subtitle/main.go inspect zh zh-Hans en-US")
	fmt.Println("  go run examples/subtitle/main.go upload BV1xx411c7mD ./demo.srt zh")
	fmt.Println("  go run examples/subtitle/main.go upload BV1xx411c7mD ./demo.srt en-US ./examples/complete/login_info.json")
	fmt.Println()
	fmt.Println("说明:")
	fmt.Println("  inspect 只打印 draft/save 接口最终会提交的 lan 值，用于离线确认映射。")
	fmt.Println("  upload 会将 SRT 转成 BCC JSON，并在提交前打印 data 摘要和 lan。")
}

func runInspect(args []string) {
	languages := args
	if len(languages) == 0 {
		languages = subtitleLanguageSamples
	}

	fmt.Println("字幕语言映射结果（提交字段: lan）:")
	for _, language := range languages {
		fmt.Printf("  %-10s -> %s\n", language, bilibili.NormalizeSubtitleLanguage(language))
	}
}

func runUpload(args []string) {
	if len(args) < 3 {
		fmt.Println("upload 模式缺少参数")
		printUsage()
		os.Exit(1)
	}

	bvid := args[0]
	subtitlePath := args[1]
	language := args[2]
	loginInfoPath := ""
	if len(args) >= 4 {
		loginInfoPath = args[3]
	}

	normalizedLanguage := bilibili.NormalizeSubtitleLanguage(language)
	fmt.Printf("原始语言值: %s\n", language)
	fmt.Printf("映射后语言值(files[].lan): %s\n", normalizedLanguage)

	loginInfo, actualLoginInfoPath, err := loadLoginInfo(loginInfoPath)
	if err != nil {
		log.Fatalf("加载登录信息失败: %v", err)
	}
	fmt.Printf("使用登录信息文件: %s\n", actualLoginInfoPath)

	client := bilibili.NewClient()
	uploader := bilibili.NewSubtitleUploader(client, loginInfo)

	videoInfo, err := uploader.GetVideoInfo(bvid)
	if err != nil {
		log.Fatalf("获取视频信息失败: %v", err)
	}

	subtitleData, err := bilibili.LoadSRTAsBCC(subtitlePath)
	if err != nil {
		log.Fatalf("解析字幕文件失败: %v", err)
	}

	fmt.Printf("目标视频: aid=%d cid=%d\n", videoInfo.AID, videoInfo.CID)
	fmt.Printf("字幕条目数: %d\n", len(subtitleData.Body))
	printSubtitleDraftSummary(subtitleData)

	acceptedLanguage, err := saveSubtitleWithCandidates(uploader, bvid, videoInfo.CID, subtitleData, language, normalizedLanguage)
	if err != nil {
		log.Fatalf("保存字幕信息失败: %v", err)
	}

	fmt.Println("字幕上传完成")
	fmt.Printf("最终被接口接受的 lan: %s\n", acceptedLanguage)
}

func printSubtitleDraftSummary(subtitle *bilibili.BCCSubtitle) {
	if subtitle == nil {
		return
	}

	summary := struct {
		FontSize float64                    `json:"font_size"`
		Count    int                        `json:"count"`
		First    *bilibili.BCCSubtitleItem `json:"first,omitempty"`
	}{
		FontSize: subtitle.FontSize,
		Count:    len(subtitle.Body),
	}
	if len(subtitle.Body) > 0 {
		summary.First = &subtitle.Body[0]
	}

	encoded, err := json.Marshal(summary)
	if err != nil {
		log.Fatalf("构造 data 摘要失败: %v", err)
	}

	fmt.Printf("即将提交的 data 摘要: %s\n", string(encoded))
}

func loadLoginInfo(path string) (*bilibili.LoginInfo, string, error) {
	resolvedPath, err := resolveLoginInfoPath(path)
	if err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, "", fmt.Errorf("read login info: %w", err)
	}

	var loginInfo bilibili.LoginInfo
	if err := json.Unmarshal(data, &loginInfo); err != nil {
		return nil, "", fmt.Errorf("unmarshal login info: %w", err)
	}

	return &loginInfo, resolvedPath, nil
}

func resolveLoginInfoPath(path string) (string, error) {
	if path != "" {
		if _, err := os.Stat(path); err != nil {
			return "", fmt.Errorf("stat login info %s: %w", path, err)
		}
		return path, nil
	}

	for _, candidate := range defaultLoginInfoCandidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("未找到 %s，可显式传入路径，或先运行 go run examples/login/main.go 生成登录信息", defaultLoginInfoPath)
}

func saveSubtitleWithCandidates(uploader *bilibili.SubtitleUploader, bvid string, cid int64, subtitle *bilibili.BCCSubtitle, originalLanguage, normalizedLanguage string) (string, error) {
	candidates := buildLanguageCandidates(originalLanguage, normalizedLanguage)
	var lastErr error

	for index, candidate := range candidates {
		fmt.Printf("即将提交的 lan: %s\n", candidate)
		fmt.Printf("尝试保存字幕[%d/%d]: lan=%s\n", index+1, len(candidates), candidate)

		if err := uploader.SaveSubtitleDraft(bvid, cid, subtitle, candidate); err == nil {
			return candidate, nil
		} else {
			lastErr = err
			fmt.Printf("保存失败: %v\n", err)
			if !isInvalidLanguageError(err) {
				return "", err
			}
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("没有可尝试的语言值")
	}

	return "", fmt.Errorf("所有候选语言值都被接口拒绝: %w", lastErr)
}

func buildLanguageCandidates(originalLanguage, normalizedLanguage string) []string {
	seen := make(map[string]struct{})
	var candidates []string

	add := func(language string) {
		language = strings.TrimSpace(language)
		if language == "" {
			return
		}
		if _, ok := seen[language]; ok {
			return
		}
		seen[language] = struct{}{}
		candidates = append(candidates, language)
	}

	primaryOriginal := primaryLanguageTag(originalLanguage)
	primaryNormalized := primaryLanguageTag(normalizedLanguage)

	add(originalLanguage)
	add(normalizedLanguage)
	add(primaryOriginal)
	add(primaryNormalized)

	switch strings.ToLower(primaryOriginal) {
	case "zh", "cmn":
		add("zh-SG")
		add("zh")
		add("zh-CN")
		add("zh-Hans")
		add("cmn-Hans")
	case "en":
		add("en")
		add("en-US")
	}

	return candidates
}

func primaryLanguageTag(language string) string {
	parts := strings.Split(strings.TrimSpace(language), "-")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func isInvalidLanguageError(err error) bool {
	if err == nil {
		return false
	}

	message := err.Error()
	return strings.Contains(message, "code=79011") || strings.Contains(message, "不合法的语言")
}