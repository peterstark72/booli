# Booli Go package

A simple Go wrapper for the [Booli API](https://www.booli.se/api).

The following functions are available:
* Listings : iterates over all "/listings"
* Sold : iterates over all "/sold"
* Areas : iterates over all "/areas"

Note that the functions returns iterators that take care of pagination in the background. If you only want the first 100 result items, you need to count the results and cancel yourself.

For each function you can specify a ```booli.Query```, which is a simple
map of string values. For example:
```
q := booli.Query{"q": "nacka"}
q := booli.Query{"q": "nacka", "minListPrice": "1000000"}
q := booli.Query{"q": "nacka", "objectType": "villa"}
```

In addition, the ```ImageURL``` method returns URL to the thumbnail picture.

### Example usage
```
package main

import (
	"fmt"

	"github.com/peterstark72/booli"
)

func main() {

	query := booli.Query{"q": "Tygelsjö"}

	for p := range booli.Listings(query) {
		fmt.Println(n, p.Location.Address.StreetAddress)

		fmt.Println(p.ImageURL())

	}
}
``` 
## Command Line tool

The command line tool can be used in the following way.

To query for listing for an area name.
```
go run cmd/main.go -q tygelsjö
```
To query for new builds for an area name.
```
go run cmd/main.go -q tygelsjö -newbuilds
```
To query for sold properties for an area name.
```
go run cmd/main.go -q tygelsjö -sold
```