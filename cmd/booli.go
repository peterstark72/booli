package main

import (
	"fmt"
	"os"

	"github.com/peterstark72/booli"
)

func main() {

	if len(os.Args) < 2 {
		panic("Usage: booli <query>")
	}

	q := booli.Query{"q": os.Args[1]}
	for p := range booli.Listings(q) {
		fmt.Println(p)
	}

}
