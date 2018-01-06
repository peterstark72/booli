package booli_test

import (
	"testing"

	"github.com/peterstark72/booli"
)

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

func TestListings(t *testing.T) {

	query := booli.Query{"q": "Malmö"}

	var solds []booli.Property
	for p := range booli.GetManyListings(query) {
		solds = append(solds, p)
	}
	if len(solds) == 0 {
		t.Error("No listings")
	}
}

func TestAreas(t *testing.T) {

	query := booli.Query{"q": "Malmö"}

	var areas []booli.Area
	for a := range booli.GetManyAreas(query) {
		areas = append(areas, a)
	}
	if len(areas) == 0 {
		t.Error("No areas")
	}
}
