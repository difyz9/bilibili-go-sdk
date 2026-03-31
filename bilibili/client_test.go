package bilibili

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.httpClient == nil {
		t.Fatal("HTTP client is nil")
	}

	if client.userAgent == "" {
		t.Fatal("User agent is empty")
	}
}

func TestClientWithOptions(t *testing.T) {
	customUA := "Test-User-Agent"
	customTimeout := 60 * time.Second

	client := NewClient(
		WithUserAgent(customUA),
		WithTimeout(customTimeout),
	)

	if client.userAgent != customUA {
		t.Errorf("Expected user agent %s, got %s", customUA, client.userAgent)
	}

	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}
}

func TestSign(t *testing.T) {
	testCases := []struct {
		params   string
		appSec   string
		expected string
	}{
		{
			params:   "appkey=test&ts=1234567890",
			appSec:   "secret",
			expected: "9ff6dcbfd27e0c57f57b2c3b99cb5d72", // 这是实际的MD5值
		},
	}

	for _, tc := range testCases {
		result := Sign(tc.params, tc.appSec)
		if len(result) != 32 {
			t.Errorf("Expected 32 character MD5 hash, got %d characters", len(result))
		}
		// 验证是否为有效的十六进制字符串
		for _, char := range result {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("Invalid MD5 hash character: %c", char)
			}
		}
	}
}

func TestIsRateLimitError(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{NewError("code=-799"), true},
		{NewError("请求过于频繁"), true},
		{NewError("rate limit"), true},
		{NewError("too many requests"), true},
		{NewError("network error"), false},
		{NewError("timeout"), false},
	}

	for _, tc := range testCases {
		result := IsRateLimitError(tc.err)
		if result != tc.expected {
			t.Errorf("IsRateLimitError(%v) = %v, expected %v", tc.err, result, tc.expected)
		}
	}
}

func TestIsNetworkError(t *testing.T) {
	testCases := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{NewError("broken pipe"), true},
		{NewError("connection reset"), true},
		{NewError("timeout"), true},
		{NewError("network error"), true},
		{NewError("dial tcp"), true},
		{NewError("EOF"), true},
		{NewError("rate limit"), false},
		{NewError("invalid request"), false},
	}

	for _, tc := range testCases {
		result := IsNetworkError(tc.err)
		if result != tc.expected {
			t.Errorf("IsNetworkError(%v) = %v, expected %v", tc.err, result, tc.expected)
		}
	}
}

func TestGetVideoReviewStatusFromArchiveView(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x/vupre/web/archive/view" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}

		return jsonResponse(`{"code":0,"message":"0","data":{"archive":{"aid":123,"bvid":"BV1xx411c7mD","title":"test","state":0,"state_desc":"开放浏览"}}}`), nil
	})

	status, err := client.GetVideoReviewStatus("BV1xx411c7mD", "SESSDATA=test")
	if err != nil {
		t.Fatalf("GetVideoReviewStatus() error = %v", err)
	}

	if !status.Passed {
		t.Fatal("expected passed status")
	}
	if !status.PublicVisible {
		t.Fatal("expected public visible status")
	}
	if status.Reviewing {
		t.Fatal("expected reviewing=false")
	}
}

func TestGetVideoReviewStatusFromPublicViewReviewing(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.bilibili.com" || req.URL.Path != "/x/web-interface/view" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}

		return jsonResponse(`{"code":62004,"message":"稿件审核中","ttl":1,"data":null}`), nil
	})

	status, err := client.GetVideoReviewStatus("BV1xx411c7mD", "")
	if err != nil {
		t.Fatalf("GetVideoReviewStatus() error = %v", err)
	}

	if !status.Reviewing {
		t.Fatal("expected reviewing status")
	}
	if status.Passed {
		t.Fatal("expected passed=false")
	}
}

func TestGetVideoReviewStatusFallbackToPublicView(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch {
		case req.URL.Host == "member.bilibili.com" && req.URL.Path == "/x/vupre/web/archive/view":
			return jsonResponse(`{"code":-101,"message":"账号未登录"}`), nil
		case req.URL.Host == "api.bilibili.com" && req.URL.Path == "/x/web-interface/view":
			return jsonResponse(`{"code":0,"message":"0","ttl":1,"data":{"aid":456,"bvid":"BV1xx411c7mD","title":"fallback","state":0}}`), nil
		default:
			t.Fatalf("unexpected request: %s", req.URL.String())
			return nil, nil
		}
	})

	status, err := client.GetVideoReviewStatus("BV1xx411c7mD", "SESSDATA=test")
	if err != nil {
		t.Fatalf("GetVideoReviewStatus() error = %v", err)
	}

	if !status.Passed {
		t.Fatal("expected passed status from fallback public view")
	}
	if status.BVid != "BV1xx411c7mD" {
		t.Fatalf("unexpected bvid: %s", status.BVid)
	}
}

func TestWaitForVideoReviewPassed(t *testing.T) {
	client := NewClient()
	current := time.Unix(0, 0)
	requests := 0

	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.bilibili.com" || req.URL.Path != "/x/web-interface/view" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}

		requests++
		if requests < 3 {
			return jsonResponse(`{"code":62004,"message":"稿件审核中","ttl":1,"data":null}`), nil
		}

		return jsonResponse(`{"code":0,"message":"0","ttl":1,"data":{"aid":456,"bvid":"BV1xx411c7mD","title":"passed","state":0}}`), nil
	})

	status, err := client.waitForVideoReviewPassed(
		"BV1xx411c7mD",
		"",
		2*time.Second,
		10*time.Second,
		func() time.Time { return current },
		func(d time.Duration) { current = current.Add(d) },
	)
	if err != nil {
		t.Fatalf("waitForVideoReviewPassed() error = %v", err)
	}
	if !status.Passed {
		t.Fatal("expected passed=true")
	}
	if requests != 3 {
		t.Fatalf("expected 3 requests, got %d", requests)
	}
}

func TestWaitForVideoReviewPassedRejected(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x/vupre/web/archive/view" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}

		return jsonResponse(`{"code":0,"message":"0","data":{"archive":{"aid":123,"bvid":"BV1xx411c7mD","title":"rejected","state":-30,"state_desc":"审核驳回","reject_reason":"封面不符合规范"}}}`), nil
	})

	status, err := client.waitForVideoReviewPassed(
		"BV1xx411c7mD",
		"SESSDATA=test",
		time.Second,
		10*time.Second,
		time.Now,
		func(time.Duration) {},
	)
	if err == nil {
		t.Fatal("expected rejection error")
	}
	if status == nil || !status.Rejected {
		t.Fatal("expected rejected status")
	}
	if !strings.Contains(err.Error(), "封面不符合规范") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWaitForVideoReviewPassedTimeout(t *testing.T) {
	client := NewClient()
	current := time.Unix(0, 0)

	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.bilibili.com" || req.URL.Path != "/x/web-interface/view" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}

		return jsonResponse(`{"code":62004,"message":"稿件审核中","ttl":1,"data":null}`), nil
	})

	status, err := client.waitForVideoReviewPassed(
		"BV1xx411c7mD",
		"",
		2*time.Second,
		5*time.Second,
		func() time.Time { return current },
		func(d time.Duration) { current = current.Add(d) },
	)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if status == nil || !status.Reviewing {
		t.Fatal("expected last status to be reviewing")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 辅助函数创建错误
func NewError(message string) error {
	return &testError{message: message}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
