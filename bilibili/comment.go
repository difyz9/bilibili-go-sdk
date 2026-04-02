package bilibili

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CommentAddRequest 发表评论请求。
type CommentAddRequest struct {
	Type    int
	OID     int64
	Root    int64
	Parent  int64
	Message string
	Plat    int
}

// CommentAddResponse 发表评论响应。
type CommentAddResponse struct {
	SuccessAction int                    `json:"success_action"`
	SuccessToast  string                 `json:"success_toast"`
	NeedCaptcha   bool                   `json:"need_captcha"`
	URL           string                 `json:"url"`
	RPID          int64                  `json:"rpid"`
	RPIDStr       string                 `json:"rpid_str"`
	Dialog        int64                  `json:"dialog"`
	DialogStr     string                 `json:"dialog_str"`
	Root          int64                  `json:"root"`
	RootStr       string                 `json:"root_str"`
	Parent        int64                  `json:"parent"`
	ParentStr     string                 `json:"parent_str"`
	Emote         map[string]interface{} `json:"emote,omitempty"`
	Reply         map[string]interface{} `json:"reply,omitempty"`
}

// CommentActionRequest 评论点赞、点踩、置顶请求。
type CommentActionRequest struct {
	Type   int
	OID    int64
	RPID   int64
	Action int
}

// CommentDeleteRequest 删除评论请求。
type CommentDeleteRequest struct {
	Type int
	OID  int64
	RPID int64
}

// CommentReportRequest 举报评论请求。
type CommentReportRequest struct {
	Type    int
	OID     int64
	RPID    int64
	Reason  int
	Content string
}

type commentAPIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    T      `json:"data"`
}

// AddComment 发表评论。
func (c *Client) AddComment(commentReq *CommentAddRequest, cookies string) (*CommentAddResponse, error) {
	if commentReq == nil {
		return nil, fmt.Errorf("comment request is required")
	}
	if commentReq.Type <= 0 {
		return nil, fmt.Errorf("comment type must be greater than 0")
	}
	if commentReq.OID <= 0 {
		return nil, fmt.Errorf("comment oid must be greater than 0")
	}
	if strings.TrimSpace(commentReq.Message) == "" {
		return nil, fmt.Errorf("comment message is required")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	formData.Set("type", strconv.Itoa(commentReq.Type))
	formData.Set("oid", strconv.FormatInt(commentReq.OID, 10))
	if commentReq.Root > 0 {
		formData.Set("root", strconv.FormatInt(commentReq.Root, 10))
	}
	if commentReq.Parent > 0 {
		formData.Set("parent", strconv.FormatInt(commentReq.Parent, 10))
	}
	formData.Set("message", commentReq.Message)
	if commentReq.Plat <= 0 {
		formData.Set("plat", "1")
	} else {
		formData.Set("plat", strconv.Itoa(commentReq.Plat))
	}
	formData.Set("csrf", csrf)

	return postCommentForm[CommentAddResponse](c, "https://api.bilibili.com/x/v2/reply/add", formData, cookies, "add comment")
}

// LikeComment 点赞或取消点赞评论。
func (c *Client) LikeComment(actionReq *CommentActionRequest, cookies string) error {
	return c.postCommentAction("https://api.bilibili.com/x/v2/reply/action", actionReq, cookies, "like comment")
}

// HateComment 点踩或取消点踩评论。
func (c *Client) HateComment(actionReq *CommentActionRequest, cookies string) error {
	return c.postCommentAction("https://api.bilibili.com/x/v2/reply/hate", actionReq, cookies, "hate comment")
}

// DeleteComment 删除评论。
func (c *Client) DeleteComment(deleteReq *CommentDeleteRequest, cookies string) error {
	if deleteReq == nil {
		return fmt.Errorf("delete request is required")
	}
	if deleteReq.Type <= 0 {
		return fmt.Errorf("comment type must be greater than 0")
	}
	if deleteReq.OID <= 0 {
		return fmt.Errorf("comment oid must be greater than 0")
	}
	if deleteReq.RPID <= 0 {
		return fmt.Errorf("comment rpid must be greater than 0")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return err
	}

	formData := url.Values{}
	formData.Set("type", strconv.Itoa(deleteReq.Type))
	formData.Set("oid", strconv.FormatInt(deleteReq.OID, 10))
	formData.Set("rpid", strconv.FormatInt(deleteReq.RPID, 10))
	formData.Set("csrf", csrf)

	_, err = postCommentForm[struct{}](c, "https://api.bilibili.com/x/v2/reply/del", formData, cookies, "delete comment")
	return err
}

// TopComment 置顶或取消置顶评论。
func (c *Client) TopComment(actionReq *CommentActionRequest, cookies string) error {
	return c.postCommentAction("https://api.bilibili.com/x/v2/reply/top", actionReq, cookies, "top comment")
}

// ReportComment 举报评论。
func (c *Client) ReportComment(reportReq *CommentReportRequest, cookies string) error {
	if reportReq == nil {
		return fmt.Errorf("report request is required")
	}
	if reportReq.Type <= 0 {
		return fmt.Errorf("comment type must be greater than 0")
	}
	if reportReq.OID <= 0 {
		return fmt.Errorf("comment oid must be greater than 0")
	}
	if reportReq.RPID <= 0 {
		return fmt.Errorf("comment rpid must be greater than 0")
	}
	if reportReq.Reason < 0 {
		return fmt.Errorf("report reason must be greater than or equal to 0")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return err
	}

	formData := url.Values{}
	formData.Set("type", strconv.Itoa(reportReq.Type))
	formData.Set("oid", strconv.FormatInt(reportReq.OID, 10))
	formData.Set("rpid", strconv.FormatInt(reportReq.RPID, 10))
	formData.Set("reason", strconv.Itoa(reportReq.Reason))
	if reportReq.Reason == 0 && strings.TrimSpace(reportReq.Content) != "" {
		formData.Set("content", reportReq.Content)
	}
	formData.Set("csrf", csrf)

	_, err = postCommentForm[struct{}](c, "https://api.bilibili.com/x/v2/reply/report", formData, cookies, "report comment")
	return err
}

func (c *Client) postCommentAction(apiURL string, actionReq *CommentActionRequest, cookies string, action string) error {
	if actionReq == nil {
		return fmt.Errorf("action request is required")
	}
	if actionReq.Type <= 0 {
		return fmt.Errorf("comment type must be greater than 0")
	}
	if actionReq.OID <= 0 {
		return fmt.Errorf("comment oid must be greater than 0")
	}
	if actionReq.RPID <= 0 {
		return fmt.Errorf("comment rpid must be greater than 0")
	}
	if actionReq.Action != 0 && actionReq.Action != 1 {
		return fmt.Errorf("action must be 0 or 1")
	}

	csrf, err := requireCSRFFromCookies(cookies)
	if err != nil {
		return err
	}

	formData := url.Values{}
	formData.Set("type", strconv.Itoa(actionReq.Type))
	formData.Set("oid", strconv.FormatInt(actionReq.OID, 10))
	formData.Set("rpid", strconv.FormatInt(actionReq.RPID, 10))
	formData.Set("action", strconv.Itoa(actionReq.Action))
	formData.Set("csrf", csrf)

	_, err = postCommentForm[struct{}](c, apiURL, formData, cookies, action)
	return err
}

func postCommentForm[T any](c *Client, apiURL string, formData url.Values, cookies string, action string) (*T, error) {
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result commentAPIResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("%s failed: code=%d, message=%s", action, result.Code, result.Message)
	}

	return &result.Data, nil
}
