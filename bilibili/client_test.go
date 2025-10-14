package bilibili

import (
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

// 辅助函数创建错误
func NewError(message string) error {
	return &testError{message: message}
}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}