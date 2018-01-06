# Booli Go package

A simple Go wrapper for the [Booli API](https://www.booli.se/api).

The following functions are available:
* GetManyListings : iterates over all "/listings"
* GetManySold : iterates over all "/sold"
* GetManyAreas : iterates over all "/areas"

Note that the functions returns iterators that take care of pagination in the background. If you only want the first 100 result items, you need to count the results and cancel yourself (see example below).

For each function you can specify a ```booli.Query```, which is a simple
map of string values. For example:
```
q := booli.Query{"q": "nacka"}
q := booli.Query{"q": "nacka", "minListPrice": "1000000"}
q := booli.Query{"q": "nacka", "objectType": "villa"}
```

In addition, the ```GetPictureURL``` function returns URL to the thumbnail picture.

### Example usage
```
package main

import (
	"fmt"

	"github.com/peterstark72/booli"
)

const MaxResults = 100

func main() {

	query := booli.Query{"q": "Malmö"}

	var n int
	for p := range booli.GetManyListings(query) {
		fmt.Println(n, p.Location.Address.StreetAddress)

		u := booli.GetPictureURL(p.BooliID)
		fmt.Println(u)

		n++
		if n > MaxResults {
			break
		}
	}
}
``` 

### Use with Google Appengine (GAE)

If you are using this client with GAE you must use the GAE HTTP service instead of the standard Go HTTP. 
You do this by creating a booli.Client and set the Transport field to the GAE HTTP client. 

```
c := appengine.NewContext(r)
client := urlfetch.Client(c)

api := booli.Client{
		Transport: client,
	}

api.GetManySold(booli.Query{"q": "nacka"})
```