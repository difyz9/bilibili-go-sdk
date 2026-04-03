package main

import (
	"fmt"
	"log"

	"github.com/difyz9/bilibili-go-sdk/bilibili"
)

func main() {

	zones := bilibili.GetVideoZones()
	fmt.Printf("main zones: %d\n", len(zones))
	for _, zone := range zones {
		fmt.Printf("- %s (tid=%d)\n", zone.Name, zone.ID)
	}

	allZones := bilibili.GetAllVideoZones()
	fmt.Printf("all zones: %d\n", len(allZones))

	zone, ok := bilibili.GetVideoZoneByTID(122)
	if !ok {
		log.Fatal("expected tid 122 to exist")
	}
	fmt.Printf("zone for tid 122: %s (parent tid=%d)\n", zone.Name, zone.ParentID)

	matchesByTID := bilibili.GetChildVideoZones(36)
	fmt.Printf("matches for tid 36: %d\n", len(matchesByTID))
	for _, match := range matchesByTID {
		fmt.Printf("-> %s (tid=%d)\n", match.Name, match.ID)
	}

	matchesBySlug := bilibili.FindVideoZonesBySlug("information")
	fmt.Printf("matches for slug 'information': %d\n", len(matchesBySlug))
	for _, match := range matchesBySlug {
		fmt.Printf("- %s (tid=%d)\n", match.Name, match.ID)
	}

}
