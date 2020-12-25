package booli_test

import (
	"testing"
	"time"

	"github.com/peterstark72/booli"
)

func TestSold(t *testing.T) {

	query := booli.Query{"q": "Tygelsjö", "objectType": "Villa"}

	var solds []booli.Property
	for p := range booli.Sold(query) {
		t.Logf("%s, %s, %s", p.Location.Address, p.ObjectType, time.Time(p.SoldDate).Format("2006-01-02"))
		solds = append(solds, p)
	}
	if len(solds) == 0 {
		t.Error("No solds")
	}
}

func TestListings(t *testing.T) {

	query := booli.Query{"q": "Tygelsjö"}

	var solds []booli.Property
	for p := range booli.Listings(query) {
		t.Logf("%s, %s, %s", p.Location.Address, p.ObjectType, time.Time(p.Published).Format("2006-01-02"))
		solds = append(solds, p)
	}
	if len(solds) == 0 {
		t.Error("No listings")
	}
}

func TestListingsAreaId(t *testing.T) {

	query := booli.Query{"areaId": "117099,866229"}

	var solds []booli.Property
	for p := range booli.Listings(query) {
		solds = append(solds, p)
	}
	if len(solds) == 0 {
		t.Error("No listings")
	}
}

func TestAreas(t *testing.T) {

	query := booli.Query{"q": "Malmö"}

	var areas []booli.Area
	for a := range booli.Areas(query) {
		areas = append(areas, a)
	}
	if len(areas) == 0 {
		t.Error("No areas")
	}
}
