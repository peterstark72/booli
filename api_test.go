package booli_test

import (
	"testing"

	"github.com/peterstark72/booli"
)

func TestListings(t *testing.T) {

	query := booli.Query{"q": "Malmö"}

	var listings []booli.Property
	for p := range booli.GetManyListings(query) {
		listings = append(listings, p)
	}
	if len(listings) == 0 {
		t.Error("No listings")
	}
}

func TestSold(t *testing.T) {

	query := booli.Query{"q": "Malmö", "minSoldDate": "20170101"}

	var solds []booli.Property
	for p := range booli.GetManySold(query) {
		solds = append(solds, p)
	}
	if len(solds) == 0 {
		t.Error("No solds")
	}
}

func TestAreas(t *testing.T) {

	query := booli.Query{"q": "Malmö"}

	var areas []booli.Area
	for p := range booli.GetManyAreas(query) {
		areas = append(areas, p)
	}
	if len(areas) == 0 {
		t.Error("No areas")
	}
}
