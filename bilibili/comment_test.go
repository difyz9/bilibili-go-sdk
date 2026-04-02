package bilibili

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestAddComment(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		if req.URL.Host != "api.bilibili.com" || req.URL.Path != "/x/v2/reply/add" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}
		if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Fatalf("unexpected content type: %s", req.Header.Get("Content-Type"))
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body failed: %v", err)
		}
		values, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		if values.Get("type") != "1" || values.Get("oid") != "243322853" || values.Get("message") != "测试评论" || values.Get("plat") != "1" || values.Get("csrf") != "csrf-token" {
			t.Fatalf("unexpected form values: %v", values)
		}

		return jsonResponse(`{"code":0,"message":"0","ttl":1,"data":{"success_action":0,"success_toast":"发送成功","need_captcha":false,"url":"","rpid":3039053308,"rpid_str":"3039053308","dialog":0,"dialog_str":"0","root":0,"root_str":"0","parent":0,"parent_str":"0"}}`), nil
	})

	resp, err := client.AddComment(&CommentAddRequest{
		Type:    1,
		OID:     243322853,
		Message: "测试评论",
	}, "SESSDATA=test; bili_jct=csrf-token")
	if err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}
	if resp.RPID != 3039053308 || resp.SuccessToast != "发送成功" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestCommentActions(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		call   func(*Client) error
		action string
	}{
		{
			name:   "like",
			path:   "/x/v2/reply/action",
			action: "1",
			call: func(client *Client) error {
				return client.LikeComment(&CommentActionRequest{Type: 1, OID: 243322853, RPID: 3039053308, Action: 1}, "SESSDATA=test; bili_jct=csrf-token")
			},
		},
		{
			name:   "hate",
			path:   "/x/v2/reply/hate",
			action: "0",
			call: func(client *Client) error {
				return client.HateComment(&CommentActionRequest{Type: 1, OID: 243322853, RPID: 3039053308, Action: 0}, "SESSDATA=test; bili_jct=csrf-token")
			},
		},
		{
			name:   "top",
			path:   "/x/v2/reply/top",
			action: "1",
			call: func(client *Client) error {
				return client.TopComment(&CommentActionRequest{Type: 1, OID: 243322853, RPID: 2940645593, Action: 1}, "SESSDATA=test; bili_jct=csrf-token")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient()
			client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.Host != "api.bilibili.com" || req.URL.Path != tc.path {
					t.Fatalf("unexpected request: %s", req.URL.String())
				}

				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("read body failed: %v", err)
				}
				values, err := url.ParseQuery(string(bodyBytes))
				if err != nil {
					t.Fatalf("parse form failed: %v", err)
				}
				if values.Get("type") != "1" || values.Get("oid") != "243322853" || values.Get("action") != tc.action || values.Get("csrf") != "csrf-token" {
					t.Fatalf("unexpected form values: %v", values)
				}

				return jsonResponse(`{"code":0,"message":"0","ttl":1}`), nil
			})

			if err := tc.call(client); err != nil {
				t.Fatalf("%s action error = %v", tc.name, err)
			}
		})
	}
}

func TestDeleteAndReportComment(t *testing.T) {
	client := NewClient()
	requestCount := 0
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requestCount++
		switch req.URL.Path {
		case "/x/v2/reply/del":
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read delete body failed: %v", err)
			}
			values, err := url.ParseQuery(string(bodyBytes))
			if err != nil {
				t.Fatalf("parse delete form failed: %v", err)
			}
			if values.Get("rpid") != "3039053308" || values.Get("csrf") != "csrf-token" {
				t.Fatalf("unexpected delete values: %v", values)
			}
			return jsonResponse(`{"code":0,"message":"0","ttl":1}`), nil
		case "/x/v2/reply/report":
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read report body failed: %v", err)
			}
			values, err := url.ParseQuery(string(bodyBytes))
			if err != nil {
				t.Fatalf("parse report form failed: %v", err)
			}
			if values.Get("reason") != "0" || values.Get("content") != "其他原因" || values.Get("csrf") != "csrf-token" {
				t.Fatalf("unexpected report values: %v", values)
			}
			return jsonResponse(`{"code":0,"message":"0","ttl":1}`), nil
		default:
			t.Fatalf("unexpected request: %s", req.URL.String())
			return nil, nil
		}
	})

	if err := client.DeleteComment(&CommentDeleteRequest{Type: 1, OID: 243322853, RPID: 3039053308}, "SESSDATA=test; bili_jct=csrf-token"); err != nil {
		t.Fatalf("DeleteComment() error = %v", err)
	}
	if err := client.ReportComment(&CommentReportRequest{Type: 1, OID: 243322853, RPID: 3039053308, Reason: 0, Content: "其他原因"}, "SESSDATA=test; bili_jct=csrf-token"); err != nil {
		t.Fatalf("ReportComment() error = %v", err)
	}
	if requestCount != 2 {
		t.Fatalf("unexpected request count: %d", requestCount)
	}
}

func TestCommentRequiresCSRF(t *testing.T) {
	client := NewClient()

	_, err := client.AddComment(&CommentAddRequest{Type: 1, OID: 1, Message: "test"}, "SESSDATA=test")
	if err == nil || !strings.Contains(err.Error(), "csrf token not found") {
		t.Fatalf("expected csrf error, got %v", err)
	}
}
