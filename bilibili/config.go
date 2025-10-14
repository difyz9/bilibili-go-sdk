// Package bilibili provides a Go SDK for Bilibili API
package bilibili

import (
	"net/http"
	"time"
)

// Config SDK配置
type Config struct {
	HTTPClient *http.Client
	UserAgent  string
	Timeout    time.Duration
	ProxyURL   string
}

// Option SDK配置选项
type Option func(*Config)

// WithHTTPClient 设置自定义HTTP客户端
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithUserAgent 设置User-Agent
func WithUserAgent(ua string) Option {
	return func(c *Config) {
		c.UserAgent = ua
	}
}

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithProxy 设置代理
func WithProxy(proxyURL string) Option {
	return func(c *Config) {
		c.ProxyURL = proxyURL
	}
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "Mozilla/5.0 BiliDroid/7.80.0 (bbcallen@gmail.com) os/android model/MI 6 mobi_app/android build/7800300 channel/bili innerVer/7800310 osVer/13 network/2",
		Timeout:   30 * time.Second,
	}
}

// ApplyOptions 应用配置选项
func (c *Config) ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	
	// 应用超时时间到HTTP客户端
	if c.HTTPClient != nil {
		c.HTTPClient.Timeout = c.Timeout
	}
}