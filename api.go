/*Package booli is a Go wrapper for the Booli API.




 */
package booli

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

//Available API resources
const (
	ListingsResource = "/listings"
	SoldResource     = "/sold"
	AreasResource    = "/areas"
)

//TimeLayout is used in response data
const TimeLayout = "2006-01-02 15:04:06"

//DateLayout is used in query parameters
const DateLayout = "20060120"

const RootURL = "https://api.booli.se"

//MaxLimitResponseSize is default size of API responses
const MaxLimitResponseSize = 100

//CallerID and PrivateKey are set in init()
var callerID, privateKey string

//letters are used to create random strings
const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	if callerID = os.Getenv("BOOLI_CALLER_ID"); callerID == "" {
		panic("Missing Booli callerID")
	}
	if privateKey = os.Getenv("BOOLI_PRIVATE_KEY"); privateKey == "" {
		panic("Missing Booli privateKey")
	}
}

type Client struct {
	Transport *http.Client
}

var DefaultClient = &Client{}

func (c Client) transport() *http.Client {
	if c.Transport != nil {
		return c.Transport
	}
	return new(http.Client)
}

//srand returns random string of size
func srand(size int) string {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte(letters[rand.Intn(len(letters))])
	}

	return string(buf)
}

//Query is Booli query parameters
type Query map[string]string

//Position in WGS84 coordinates
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Address struct {
	StreetAddress string `json:"streetAddress"`
}

//Location is collection of location elements
type Location struct {
	Position   Position `json:"position"`
	NamedAreas []string `json:"namedAreas"`
	Address    Address  `json:"address"`
}

//Property is an actual property
type Property struct {
	Location  Location `json:"location"`
	ListPrice int      `json:"listPrice"`
	SoldPrice int      `json:"soldPrice"`
	SoldDate  string   `json:"soldDate"`
	BooliID   int      `json:"booliId"`
	Published string   `json:"published"`
	URL       string   `json:"url"`
}

//Area is an area
type Area struct {
	BooliID       int      `json:"booliId"`
	Name          string   `json:"string"`
	Types         []string `json:"types"`
	ParentBooliID int      `json:"parentBooliId"`
	ParentName    string   `json:"parentName"`
	ParentTypes   string   `json:"parentTypes"`
	FullName      string   `json:"fullName"`
}

//Pagination defines https://www.booli.se/api/#pagination
type Pagination struct {
	Count      int `json:"count"`
	TotalCount int `json:"totalCount"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
}

type Response struct {
	Pagination
	Sold     []Property `json:"sold"`
	Listings []Property `json:"listings"`
	Areas    []Area     `json:"areas"`
}

/*load gets one page of response data from path.
It returns the data and an error code.
*/
func (c Client) get(path string, params Query) ([]byte, error) {

	//Create auth values
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	unique := srand(16)
	s := callerID + timestamp + privateKey + unique
	h := fmt.Sprintf("%x", sha1.Sum([]byte(s)))

	//Auth URL query values
	q := url.Values{}
	q.Set("callerId", callerID)
	q.Set("unique", unique)
	q.Set("hash", h)
	q.Set("time", timestamp)

	//Add filter params into query
	for k, v := range params {
		q.Set(k, v)
	}

	//Build URL and make request
	u := RootURL + path + "?" + q.Encode()
	res, err := c.transport().Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)

}

//paginator returns an iterator for all page responses
func (c Client) paginator(resource string, params Query) chan Response {

	ch := make(chan Response)

	go func() {

		var offset int
		limit := MaxLimitResponseSize
		for {

			//Set pagination parameters
			params["limit"] = strconv.Itoa(limit)
			params["offset"] = strconv.Itoa(offset)

			data, err := c.get(resource, params)
			if err != nil {
				break
			}

			var resp Response
			json.Unmarshal(data, &resp)
			ch <- resp

			offset += limit
			if offset > resp.TotalCount {
				break
			}
		}
		close(ch)
	}()
	return ch
}

//GetManySold iterates sold properties
func (c Client) GetManySold(params Query) chan Property {
	ch := make(chan Property)
	go func() {
		for resp := range c.paginator(SoldResource, params) {
			for _, a := range resp.Sold {
				ch <- a
			}
		}
		close(ch)
	}()
	return ch
}

//GetManySold iterates sold properties
func GetManySold(params Query) chan Property {
	return DefaultClient.GetManySold(params)
}

//GetManyListings iterates listed properties
func (c Client) GetManyListings(params Query) chan Property {
	ch := make(chan Property)
	go func() {
		for resp := range c.paginator(ListingsResource, params) {
			for _, p := range resp.Listings {
				ch <- p
			}
		}
		close(ch)
	}()
	return ch
}

//GetManyListings iterates sold properties
func GetManyListings(params Query) chan Property {
	return DefaultClient.GetManyListings(params)
}

//GetManyAreas iterates areas
func (c Client) GetManyAreas(params Query) chan Area {
	ch := make(chan Area)
	go func() {
		for resp := range c.paginator(AreasResource, params) {
			for _, a := range resp.Areas {
				ch <- a
			}
		}
		close(ch)
	}()
	return ch
}

//GetManyListings iterates sold properties
func GetManyAreas(params Query) chan Area {
	return DefaultClient.GetManyAreas(params)
}

//GetPictureURL builds picture URL: https://www.booli.se/api/#images
func GetPictureURL(booliID int) string {
	return fmt.Sprintf("https://api.bcdn.se/cache/primary_%v_140x94.jpg", booliID)
}
