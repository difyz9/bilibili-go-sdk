package bilibili

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestGetSeasonList(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x2/creative/web/seasons" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}
		if got := req.URL.Query().Get("pn"); got != "2" {
			t.Fatalf("unexpected pn: %s", got)
		}
		if got := req.URL.Query().Get("ps"); got != "10" {
			t.Fatalf("unexpected ps: %s", got)
		}

		return jsonResponse(`{"code":0,"message":"0","data":{"seasons":[{"season":{"id":1,"title":"合集A","cover":"cover.jpg"},"checkin":{"status":0,"season_status":1},"seasonStat":{"view":123},"sections":{"sections":[{"id":11,"type":1,"seasonId":1,"title":"正片","order":1,"state":0,"partState":0,"epCount":1,"cover":"cover.jpg","has_charging_pay":0}]},"part_episodes":[{"id":21,"title":"视频A","aid":1001,"bvid":"BV1xx411c7mD","cid":2001,"seasonId":1,"sectionId":11,"order":1,"archiveState":0,"state":0,"is_free":0,"aid_owner":true,"charging_pay":0}]}],"tip":{"title":"","url":""},"total":1,"play_type":1}}`), nil
	})

	data, err := client.GetSeasonList(&SeasonListParams{Page: 2, PageSize: 10}, "SESSDATA=test")
	if err != nil {
		t.Fatalf("GetSeasonList() error = %v", err)
	}
	if data.Total != 1 {
		t.Fatalf("unexpected total: %d", data.Total)
	}
	if len(data.Seasons) != 1 || data.Seasons[0].Season.Title != "合集A" {
		t.Fatalf("unexpected seasons: %+v", data.Seasons)
	}
	if len(data.Seasons[0].Sections.Sections) != 1 {
		t.Fatalf("unexpected sections: %+v", data.Seasons[0].Sections)
	}
	if len(data.Seasons[0].PartEpisodes) != 1 {
		t.Fatalf("unexpected episodes: %+v", data.Seasons[0].PartEpisodes)
	}
}

func TestCreateSeason(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x2/creative/web/season/add" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}
		if got := req.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
			t.Fatalf("unexpected content type: %s", got)
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body failed: %v", err)
		}
		values, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		if values.Get("title") != "测试合集" || values.Get("cover") != "cover.jpg" || values.Get("csrf") != "csrf-token" {
			t.Fatalf("unexpected form values: %v", values)
		}

		return jsonResponse(`{"code":0,"message":"0","data":3541247}`), nil
	})

	seasonID, err := client.CreateSeason(&SeasonCreateRequest{
		Title:       "测试合集",
		Desc:        "合集简介",
		Cover:       "cover.jpg",
		SeasonPrice: 0,
	}, "SESSDATA=test; bili_jct=csrf-token")
	if err != nil {
		t.Fatalf("CreateSeason() error = %v", err)
	}
	if seasonID != 3541247 {
		t.Fatalf("unexpected season ID: %d", seasonID)
	}
}

func TestAddSeasonEpisodes(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		if req.URL.Host != "member.bilibili.com" || req.URL.Path != "/x2/creative/web/season/section/episodes/add" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}
		if req.URL.Query().Get("csrf") != "csrf-token" {
			t.Fatalf("unexpected csrf query: %s", req.URL.RawQuery)
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body failed: %v", err)
		}
		body := string(bodyBytes)
		for _, expected := range []string{`"sectionId":3954033`, `"section_id":3954033`, `"episodes"`, `"episode"`, `"title":"视频1"`} {
			if !strings.Contains(body, expected) {
				t.Fatalf("expected %q in body: %s", expected, body)
			}
		}

		return jsonResponse(`{"code":0,"message":"0"}`), nil
	})

	err := client.AddSeasonEpisodes(&SeasonAddEpisodesRequest{
		SectionID: 3954033,
		Episodes: []SeasonEpisodeInput{{
			Title:       "视频1",
			AID:         1906473802,
			CID:         1625992822,
			ChargingPay: 0,
		}},
	}, "SESSDATA=test; bili_jct=csrf-token")
	if err != nil {
		t.Fatalf("AddSeasonEpisodes() error = %v", err)
	}
}

func TestEditSeasonSection(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/x2/creative/web/season/section/edit" {
			t.Fatalf("unexpected request: %s", req.URL.String())
		}

		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body failed: %v", err)
		}
		body := string(bodyBytes)
		for _, expected := range []string{`"section":{`, `"id":3954033`, `"type":1`, `"seasonId":3541247`, `"title":"正片"`, `"order":1`, `"sort":1`} {
			if !strings.Contains(body, expected) {
				t.Fatalf("expected %q in body: %s", expected, body)
			}
		}

		return jsonResponse(`{"code":0,"message":"0"}`), nil
	})

	err := client.EditSeasonSection(&SeasonSectionEditRequest{
		Section: SeasonSectionEditInfo{ID: 3954033, SeasonID: 3541247, Title: "正片"},
		Sorts:   []SeasonEpisodeOrder{{ID: 77260687, Order: 1}},
	}, "SESSDATA=test; bili_jct=csrf-token")
	if err != nil {
		t.Fatalf("EditSeasonSection() error = %v", err)
	}
}

func TestEditAndDeleteSeasonAndGetSection(t *testing.T) {
	client := NewClient()
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/x2/creative/web/season/edit":
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body failed: %v", err)
			}
			body := string(bodyBytes)
			for _, expected := range []string{`"season":{`, `"id":3541327`, `"title":"测试合集"`, `"cover":"cover.jpg"`, `"desc":"简介"`, `"sorts":[{"id":3954127,"sort":1}]`} {
				if !strings.Contains(body, expected) {
					t.Fatalf("expected %q in body: %s", expected, body)
				}
			}
			return jsonResponse(`{"code":0,"message":"0"}`), nil
		case "/x2/creative/web/season/del":
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read body failed: %v", err)
			}
			values, err := url.ParseQuery(string(bodyBytes))
			if err != nil {
				t.Fatalf("parse form failed: %v", err)
			}
			if values.Get("id") != "3541327" || values.Get("csrf") != "csrf-token" {
				t.Fatalf("unexpected form values: %v", values)
			}
			return jsonResponse(`{"code":0,"message":"0"}`), nil
		case "/x2/creative/web/season/section":
			if req.URL.Query().Get("id") != "176088" {
				t.Fatalf("unexpected query: %s", req.URL.RawQuery)
			}
			return jsonResponse(`{"code":0,"message":"0","data":{"section":{"id":176088,"type":1,"seasonId":152812,"title":"正片","order":1,"state":0,"partState":0,"rejectReason":"","ctime":1643250822,"mtime":1739466002,"epCount":1,"cover":"cover.jpg","has_charging_pay":0,"show":1,"has_pugv_pay":0},"episodes":[{"id":109100674,"title":"视频A","aid":113997323963614,"bvid":"BV14BNfeSE5c","cid":28376042631,"seasonId":152812,"sectionId":176088,"order":1,"videoTitle":"视频A","archiveTitle":"视频A","archiveState":0,"rejectReason":"","state":0,"cover":"","is_free":0,"aid_owner":true,"charging_pay":0}]}}`), nil
		default:
			t.Fatalf("unexpected request: %s", req.URL.String())
			return nil, nil
		}
	})

	err := client.EditSeason(&SeasonEditRequest{
		Season: SeasonEditInfo{ID: 3541327, Title: "测试合集", Cover: "cover.jpg", Desc: "简介"},
		Sorts:  []SeasonSectionOrder{{ID: 3954127, Sort: 1}},
	}, "SESSDATA=test; bili_jct=csrf-token")
	if err != nil {
		t.Fatalf("EditSeason() error = %v", err)
	}

	err = client.DeleteSeason(3541327, "SESSDATA=test; bili_jct=csrf-token")
	if err != nil {
		t.Fatalf("DeleteSeason() error = %v", err)
	}

	detail, err := client.GetSeasonSection(176088, "SESSDATA=test")
	if err != nil {
		t.Fatalf("GetSeasonSection() error = %v", err)
	}
	if detail.Section.ID != 176088 || len(detail.Episodes) != 1 {
		t.Fatalf("unexpected detail: %+v", detail)
	}
	if detail.Episodes[0].BVID != "BV14BNfeSE5c" {
		t.Fatalf("unexpected episode: %+v", detail.Episodes[0])
	}
}
