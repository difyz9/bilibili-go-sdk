package bilibili

import (
	"testing"
)

func TestNormalizeSubtitleLanguage(t *testing.T) {
	tests := map[string]string{
		"zh":       SubtitleLangZh,
		"zh-CN":    SubtitleLangZh,
		"zh-Hans":  SubtitleLangZh,
		"cmn-Hans": SubtitleLangZh,
		"zh-TW":    SubtitleLangZhTW,
		"zh-Hant":  SubtitleLangZhTW,
		"en":       SubtitleLangEN,
		"en-US":    SubtitleLangENUS,
		"ja":       "ja",
		"":         "",
	}

	for input, expected := range tests {
		if actual := NormalizeSubtitleLanguage(input); actual != expected {
			t.Fatalf("NormalizeSubtitleLanguage(%q) = %q, want %q", input, actual, expected)
		}
	}
}

func TestParseSRTToBCC(t *testing.T) {
	content := "1\n00:00:00,000 --> 00:00:01,500\n你好\n\n2\n00:00:01,500 --> 00:00:03,000\n世界\n"

	bcc, err := ParseSRTToBCC(content)
	if err != nil {
		t.Fatalf("ParseSRTToBCC returned error: %v", err)
	}

	if len(bcc.Body) != 2 {
		t.Fatalf("expected 2 subtitle items, got %d", len(bcc.Body))
	}

	if bcc.Body[0].From != 0 || bcc.Body[0].To != 1.5 || bcc.Body[0].Content != "你好" {
		t.Fatalf("unexpected first subtitle item: %+v", bcc.Body[0])
	}

	if bcc.Body[1].From != 1.5 || bcc.Body[1].To != 3 || bcc.Body[1].Content != "世界" {
		t.Fatalf("unexpected second subtitle item: %+v", bcc.Body[1])
	}

	if bcc.FontColor != "#FFFFFF" || bcc.Stroke != "none" {
		t.Fatalf("unexpected bcc metadata: %+v", bcc)
	}
}
