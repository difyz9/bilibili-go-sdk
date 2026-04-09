package bilibili

import "testing"

func TestNormalizeSubtitleLanguage(t *testing.T) {
	tests := map[string]string{
		"zh":       SubtitleLangZhCN,
		"zh-CN":    SubtitleLangZhCN,
		"zh-Hans":  SubtitleLangZhCN,
		"cmn-Hans": SubtitleLangZhCN,
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
