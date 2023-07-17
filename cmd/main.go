package main

import (
	"fmt"
	"os"
	"time"

	"github.com/peterstark72/booli"
)

func main() {

	if len(os.Args) < 2 {
		panic("Usage: booli <query>")
	}

	q := booli.Query{"q": os.Args[1]}
	for p := range booli.Listings(q) {
		fmt.Printf("%s %s %s %d/%d %s\n", p.Location.Address.StreetAddress, p.ObjectType, time.Time(p.Published).Format("2006-01-02"), p.ListPrice, p.SoldPrice, p.URL)
	}

}
