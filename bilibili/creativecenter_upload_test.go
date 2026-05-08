package bilibili

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strings"
	"testing"
)

func TestSubmitVideoUsesWebEndpoint(t *testing.T) {
	loginInfo := &LoginInfo{
		CookieInfo: map[string]interface{}{
			"cookies": []interface{}{
				map[string]interface{}{"name": "SESSDATA", "value": "sess"},
				map[string]interface{}{"name": "bili_jct", "value": "csrf-token"},
			},
		},
	}

	uploader := NewUploadClient(loginInfo)
	uploader.client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x/vu/web/add/v3" {
			t.Fatalf("unexpected url: %s", req.URL.String())
		}
		if req.URL.Query().Get("csrf") != "csrf-token" {
			t.Fatalf("missing csrf query: %s", req.URL.RawQuery)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body failed: %v", err)
		}

		var studio Studio
		if err := json.Unmarshal(body, &studio); err != nil {
			t.Fatalf("unmarshal body failed: %v", err)
		}
		if studio.DescFormatId != 9999 {
			t.Fatalf("expected desc_format_id=9999, got %d", studio.DescFormatId)
		}
		if studio.WebOS != 3 {
			t.Fatalf("expected web_os=3, got %d", studio.WebOS)
		}
		if len(studio.Videos) != 1 || studio.Videos[0].CID != 12345 {
			t.Fatalf("unexpected videos payload: %+v", studio.Videos)
		}

		return jsonResponse(`{"code":0,"message":"0","data":{"aid":1}}`), nil
	})

	_, err := uploader.SubmitVideo(&Studio{
		Title:     "test title",
		Copyright: 1,
		Tid:       122,
		Tag:       "go",
		Subtitle:  Subtitle{Open: 0, Lan: ""},
		Videos: []Video{{
			Title:    "P1",
			Filename: "upload-file-name",
			CID:      12345,
		}},
	})
	if err != nil {
		t.Fatalf("SubmitVideo() error = %v", err)
	}
}

func TestCheckTagUsesNestedDataCode(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x/vupre/web/topic/tag/check" {
			t.Fatalf("unexpected url: %s", req.URL.String())
		}
		return jsonResponse(`{"code":0,"message":"0","data":{"code":1,"content":"invalid"}}`), nil
	})

	valid, err := client.CheckTag("bad-tag")
	if err != nil {
		t.Fatalf("CheckTag() error = %v", err)
	}
	if valid {
		t.Fatal("expected tag to be invalid")
	}
}

func TestPredictArchiveTypesUsesMultipart(t *testing.T) {
	loginInfo := &LoginInfo{
		CookieInfo: map[string]interface{}{
			"cookies": []interface{}{
				map[string]interface{}{"name": "SESSDATA", "value": "sess"},
				map[string]interface{}{"name": "bili_jct", "value": "csrf-token"},
			},
		},
	}

	uploader := NewUploadClient(loginInfo)
	uploader.client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x/vupre/web/archive/types/predict" {
			t.Fatalf("unexpected url: %s", req.URL.String())
		}
		if req.URL.Query().Get("csrf") != "csrf-token" {
			t.Fatalf("unexpected query: %s", req.URL.RawQuery)
		}

		mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			t.Fatalf("parse content type failed: %v", err)
		}
		if mediaType != "multipart/form-data" {
			t.Fatalf("unexpected content type: %s", mediaType)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body failed: %v", err)
		}
		if !strings.Contains(string(body), `name="filename"`) || !strings.Contains(string(body), `name="upload_id"`) {
			t.Fatalf("unexpected multipart body: %s", string(body))
		}
		if params["boundary"] == "" {
			t.Fatal("expected multipart boundary")
		}

		return jsonResponse(`{"code":0,"message":"0","request_id":"rid","data":[{"id":122,"parent":17,"parent_name":"科技","name":"软件应用","description":"desc"}]}`), nil
	})

	result, err := uploader.PredictArchiveTypes(&ArchiveTypePredictRequest{
		Filename: "upload-file-name",
		Title:    "Predict title",
		UploadID: "upload-id-1",
	})
	if err != nil {
		t.Fatalf("PredictArchiveTypes() error = %v", err)
	}
	if len(result) != 1 || result[0].ID != 122 {
		t.Fatalf("unexpected result: %+v", result)
	}
}