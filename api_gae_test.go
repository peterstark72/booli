package booli_test

import (
	"testing"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/urlfetch"

	"github.com/peterstark72/booli"
)

//TestSoldWithGAE tests with Google App engine context
func TestSoldWithGAE(t *testing.T) {

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	client := booli.Client{
		Transport: urlfetch.Client(ctx),
	}

	query := booli.Query{"q": "Malmö", "minSoldDate": "20170101"}

	var solds []booli.Property
	for p := range client.GetManySold(query) {
		solds = append(solds, p)
	}
	if len(solds) == 0 {
		t.Error("No GAE solds")
	}
}

//TestListingsWithGAE tests with Google App engine context
func TestListingsWithGAE(t *testing.T) {

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	client := booli.Client{
		Transport: urlfetch.Client(ctx),
	}

	query := booli.Query{"q": "Malmö"}

	var solds []booli.Property
	for p := range client.GetManyListings(query) {
		solds = append(solds, p)
	}
	if len(solds) == 0 {
		t.Error("No GAE listings")
	}
}
