package main

import (
	"encoding/json"
	"os"

	"github.com/peterstark72/booli"
)

func main() {

	q := booli.Query{"q": "Tygelsjö"}

	properties := booli.GetAllListings(q)

	d, _ := json.Marshal(properties)

	os.Stdout.Write(d)
}
