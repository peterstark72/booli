package main

import (
	"flag"
	"fmt"

	"github.com/peterstark72/booli"
)

var sold bool
var query string
var newbuilds bool

func init() {
	flag.BoolVar(&sold, "sold", false, "Solds")
	flag.BoolVar(&newbuilds, "newbuilds", false, "New-builds")
	flag.StringVar(&query, "q", "", "Query")
}

func main() {

	flag.Parse()

	q := booli.Query{
		"q": query,
	}

	var items chan booli.Property
	if sold {
		items = booli.Sold(q)
	} else {
		items = booli.Listings(q)
	}

	if newbuilds {
		q["isNewConstruction"] = "1"
	}

	for p := range items {
		fmt.Printf("%s\n", p)
	}

}
