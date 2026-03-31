package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// VideoReviewStatus 投稿审核状态
type VideoReviewStatus struct {
	AID             int64  `json:"aid"`
	BVid            string `json:"bvid"`
	Title           string `json:"title,omitempty"`
	State           int    `json:"state"`
	StateDesc       string `json:"state_desc,omitempty"`
	RejectReason    string `json:"reject_reason,omitempty"`
	Passed          bool   `json:"passed"`
	Reviewing       bool   `json:"reviewing"`
	Rejected        bool   `json:"rejected"`
	PublicVisible   bool   `json:"public_visible"`
	OnlySelfVisible bool   `json:"only_self_visible"`
}

type archiveReviewStatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Archive struct {
			AID          int64  `json:"aid"`
			BVid         string `json:"bvid"`
			Title        string `json:"title"`
			State        int    `json:"state"`
			StateDesc    string `json:"state_desc"`
			RejectReason string `json:"reject_reason"`
		} `json:"archive"`
	} `json:"data"`
}

type publicVideoReviewResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AID   int64  `json:"aid"`
		BVid  string `json:"bvid"`
		Title string `json:"title"`
		State int    `json:"state"`
	} `json:"data"`
}

// GetVideoReviewStatus 获取视频投稿是否通过审核。
//
// 优先使用创作中心稿件接口获取更准确的审核状态；当未提供 Cookie 或创作中心接口不可用时，
// 回退到公开视频信息接口，并结合文档中的 code/state 判断当前是否已通过审核。
func (c *Client) GetVideoReviewStatus(bvid string, cookies string) (*VideoReviewStatus, error) {
	bvid = strings.TrimSpace(bvid)
	if bvid == "" {
		return nil, fmt.Errorf("bvid is required")
	}

	var archiveErr error
	if strings.TrimSpace(cookies) != "" {
		status, err := c.getVideoReviewStatusFromArchiveView(bvid, cookies)
		if err == nil {
			return status, nil
		}
		archiveErr = err
	}

	status, err := c.getVideoReviewStatusFromPublicView(bvid, cookies)
	if err == nil {
		return status, nil
	}
	if archiveErr != nil {
		return nil, fmt.Errorf("get video review status failed: archive view: %v; public view: %w", archiveErr, err)
	}

	return nil, err
}

// WaitForVideoReviewPassed 轮询等待视频投稿审核通过。
//
// 当检测到已通过审核时立即返回；当检测到驳回时直接返回当前状态和错误；
// 如果在 timeout 内仍未通过，则返回最后一次查询到的状态和超时错误。
func (c *Client) WaitForVideoReviewPassed(bvid string, cookies string, interval time.Duration, timeout time.Duration) (*VideoReviewStatus, error) {
	return c.waitForVideoReviewPassed(bvid, cookies, interval, timeout, time.Now, time.Sleep)
}

func (c *Client) waitForVideoReviewPassed(
	bvid string,
	cookies string,
	interval time.Duration,
	timeout time.Duration,
	now func() time.Time,
	sleep func(time.Duration),
) (*VideoReviewStatus, error) {
	bvid = strings.TrimSpace(bvid)
	if bvid == "" {
		return nil, fmt.Errorf("bvid is required")
	}
	if interval <= 0 {
		interval = 10 * time.Second
	}
	if timeout <= 0 {
		timeout = 30 * time.Minute
	}

	deadline := now().Add(timeout)

	for {
		status, err := c.GetVideoReviewStatus(bvid, cookies)
		if err != nil {
			return nil, err
		}

		if status.Passed {
			return status, nil
		}

		if status.Rejected {
			return status, fmt.Errorf("video review rejected: %s", reviewStatusMessage(status))
		}

		if !now().Before(deadline) {
			return status, fmt.Errorf("wait for video review passed timed out after %s: %s", timeout, reviewStatusMessage(status))
		}

		sleep(interval)
	}
}

func (c *Client) getVideoReviewStatusFromArchiveView(bvid string, cookies string) (*VideoReviewStatus, error) {
	apiURL := fmt.Sprintf("https://member.bilibili.com/x/vupre/web/archive/view?bvid=%s", bvid)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://member.bilibili.com/")
	req.Header.Set("Cookie", cookies)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result archiveReviewStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("get archive review status failed: code=%d, message=%s", result.Code, result.Message)
	}

	status := &VideoReviewStatus{
		AID:          result.Data.Archive.AID,
		BVid:         result.Data.Archive.BVid,
		Title:        result.Data.Archive.Title,
		State:        result.Data.Archive.State,
		StateDesc:    strings.TrimSpace(result.Data.Archive.StateDesc),
		RejectReason: strings.TrimSpace(result.Data.Archive.RejectReason),
	}

	normalizeVideoReviewStatus(status)
	return status, nil
}

func (c *Client) getVideoReviewStatusFromPublicView(bvid string, cookies string) (*VideoReviewStatus, error) {
	apiURL := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvid)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	if strings.TrimSpace(cookies) != "" {
		req.Header.Set("Cookie", cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result publicVideoReviewResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	switch result.Code {
	case 0:
		status := &VideoReviewStatus{
			AID:           result.Data.AID,
			BVid:          result.Data.BVid,
			Title:         result.Data.Title,
			State:         result.Data.State,
			StateDesc:     publicStateDescription(result.Data.State),
			PublicVisible: result.Data.State == 0,
		}
		normalizeVideoReviewStatus(status)
		return status, nil
	case 62004:
		status := &VideoReviewStatus{
			BVid:      bvid,
			StateDesc: "审核中",
		}
		normalizeVideoReviewStatus(status)
		return status, nil
	case 62012:
		status := &VideoReviewStatus{
			BVid:            bvid,
			StateDesc:       "仅UP主自己可见",
			OnlySelfVisible: true,
		}
		normalizeVideoReviewStatus(status)
		return status, nil
	default:
		return nil, fmt.Errorf("get public video info failed: code=%d, message=%s", result.Code, result.Message)
	}
}

func normalizeVideoReviewStatus(status *VideoReviewStatus) {
	desc := strings.TrimSpace(status.StateDesc)

	if containsAny(desc, "审核中", "转码中", "处理中", "排队中") {
		status.Reviewing = true
	}

	if containsAny(desc, "驳回", "未通过", "退回", "打回", "锁定") || strings.TrimSpace(status.RejectReason) != "" {
		status.Rejected = true
	}

	if containsAny(desc, "仅UP主自己可见", "仅up主自己可见", "仅自己可见") {
		status.OnlySelfVisible = true
	}

	if status.State == 0 || containsAny(desc, "开放浏览", "审核通过", "已通过") {
		status.Passed = !status.Reviewing && !status.Rejected && !status.OnlySelfVisible
	}

	if status.Passed {
		status.PublicVisible = true
	}
}

func containsAny(text string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func publicStateDescription(state int) string {
	if state == 0 {
		return "开放浏览"
	}

	return fmt.Sprintf("state=%d", state)
}

func reviewStatusMessage(status *VideoReviewStatus) string {
	if status == nil {
		return "unknown status"
	}
	if strings.TrimSpace(status.RejectReason) != "" {
		return strings.TrimSpace(status.RejectReason)
	}
	if strings.TrimSpace(status.StateDesc) != "" {
		return strings.TrimSpace(status.StateDesc)
	}
	if status.State != 0 {
		return fmt.Sprintf("state=%d", status.State)
	}
	return "unknown status"
}
