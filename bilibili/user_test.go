package bilibili

import (
	"net/http"
	"testing"
	"time"
)

func TestGetUserBasicInfo(t *testing.T) {
	client := NewClient()
	if err := client.wbiManager.UpdateKeys(
		"https://i0.hdslb.com/bfs/wbi/abcdefghijklmnopqrstuvwxyz123456.png",
		"https://i0.hdslb.com/bfs/wbi/123456abcdefghijklmnopqrstuvwxyz.png",
	); err != nil {
		t.Fatalf("UpdateKeys() error = %v", err)
	}
	client.wbiManager.lastUpdate = time.Now()

	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		if req.URL.Host != "api.bilibili.com" || req.URL.Path != "/x/space/wbi/acc/info" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}
		if got := req.URL.Query().Get("mid"); got != "2" {
			t.Fatalf("unexpected mid: %s", got)
		}
		if req.URL.Query().Get("w_rid") == "" {
			t.Fatal("expected w_rid query param")
		}
		if req.URL.Query().Get("wts") == "" {
			t.Fatal("expected wts query param")
		}
		if got := req.Header.Get("Cookie"); got != "SESSDATA=test" {
			t.Fatalf("unexpected cookie: %s", got)
		}

		return jsonResponse(`{"code":0,"message":"0","ttl":1,"data":{"mid":2,"name":"test-user","sex":"保密","face":"https://i0.hdslb.com/test.jpg","sign":"hello","rank":10000,"level":6,"birthday":"01-02","is_followed":true,"official":{"role":1,"title":"UP主认证","desc":"认证说明","type":0},"vip":{"type":2,"status":1,"due_date":1735689600000,"label":{"text":"年度大会员","label_theme":"annual_vip"}},"live_room":{"roomStatus":1,"liveStatus":0,"url":"https://live.bilibili.com/1","title":"直播间标题","roomid":1,"roundStatus":0,"broadcast_type":0}}}`), nil
	})

	info, err := client.GetUserBasicInfo(2, "SESSDATA=test")
	if err != nil {
		t.Fatalf("GetUserBasicInfo() error = %v", err)
	}
	if info.Mid != 2 || info.Name != "test-user" {
		t.Fatalf("unexpected user info: %+v", info)
	}
	if info.VIP == nil || info.VIP.Label == nil || info.VIP.Label.Text != "年度大会员" {
		t.Fatalf("unexpected vip info: %+v", info.VIP)
	}
	if info.LiveRoom == nil || info.LiveRoom.RoomID != 1 {
		t.Fatalf("unexpected live room info: %+v", info.LiveRoom)
	}
}

func TestGetUserBasicInfoInvalidMID(t *testing.T) {
	client := NewClient()

	info, err := client.GetUserBasicInfo(0, "SESSDATA=test")
	if err == nil {
		t.Fatal("expected error for invalid mid")
	}
	if info != nil {
		t.Fatalf("expected nil info, got %+v", info)
	}
}
