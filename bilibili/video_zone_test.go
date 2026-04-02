package bilibili

import "testing"

func TestGetVideoZones(t *testing.T) {
	zones := GetVideoZones()
	if len(zones) != 22 {
		t.Fatalf("unexpected main zone count: %d", len(zones))
	}

	allZones := GetAllVideoZones()
	if len(allZones) != 169 {
		t.Fatalf("unexpected flattened zone count: %d", len(allZones))
	}

	if len(zones[0].Children) != 8 {
		t.Fatalf("unexpected child zone count for 动画: %d", len(zones[0].Children))
	}

	zones[0].Name = "mutated"
	zones[0].Children[0].Name = "mutated child"
	fresh := GetVideoZones()
	if fresh[0].Name != "动画" {
		t.Fatalf("expected cloned main zone data, got %q", fresh[0].Name)
	}
	if fresh[0].Children[0].Name != "MAD·AMV" {
		t.Fatalf("expected cloned child zone data, got %q", fresh[0].Children[0].Name)
	}
}

func TestGetVideoZoneByTID(t *testing.T) {
	zone, ok := GetVideoZoneByTID(265)
	if !ok {
		t.Fatal("expected tid 265 to exist")
	}
	if zone.Name != "AI音乐" || zone.ParentID != 3 {
		t.Fatalf("unexpected zone: %+v", zone)
	}

	zone176, ok := GetVideoZoneByTID(176)
	if !ok {
		t.Fatal("expected tid 176 to exist")
	}
	if zone176.Name != "汽车生活" || zone176.IsRedirect {
		t.Fatalf("expected best match for tid 176 to prefer active zone, got %+v", zone176)
	}

	if _, ok := GetVideoZoneByTID(999999); ok {
		t.Fatal("expected missing tid lookup to fail")
	}
}

func TestFindVideoZones(t *testing.T) {
	matchesByTID := FindVideoZonesByTID(176)
	if len(matchesByTID) != 2 {
		t.Fatalf("expected 2 matches for tid 176, got %d", len(matchesByTID))
	}

	matchesBySlug := FindVideoZonesBySlug("information")
	if len(matchesBySlug) != 3 {
		t.Fatalf("expected 3 matches for slug information, got %d", len(matchesBySlug))
	}

	matchesByName := FindVideoZonesByName("资讯")
	if len(matchesByName) != 3 {
		t.Fatalf("expected 3 matches for name 资讯, got %d", len(matchesByName))
	}

	matchesByURI := FindVideoZonesByURI("/v/sports/culture")
	if len(matchesByURI) != 2 {
		t.Fatalf("expected 2 matches for /v/sports/culture, got %d", len(matchesByURI))
	}

	children := GetChildVideoZones(1)
	if len(children) != 8 {
		t.Fatalf("unexpected child count for tid 1: %d", len(children))
	}

	deprecated := FindVideoZonesByTID(194)
	if len(deprecated) != 1 || !deprecated[0].IsDeprecated {
		t.Fatalf("expected tid 194 to be deprecated, got %+v", deprecated)
	}

	redirect := FindVideoZonesByTID(163)
	if len(redirect) != 1 || !redirect[0].IsRedirect {
		t.Fatalf("expected tid 163 to be redirect, got %+v", redirect)
	}
}
