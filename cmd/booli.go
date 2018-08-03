package main

import (
	"fmt"

	"github.com/peterstark72/booli"
)

func main() {

	q := booli.Query{"q": "Tygelsjö"}

	for _, obj := range booli.GetAllListings(q) {
		fmt.Println(obj)
		fmt.Printf("https://api.bcdn.se/cache/primary_%v_140x94.jpg\n", obj.BooliID)
	}

}
