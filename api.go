// Package booli is a Go wrapper for the Booli API.
package booli

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Available API resources
const (
	ListingsResource = "listings"
	SoldResource     = "sold"
	AreasResource    = "areas"
)

// RootURL is Booli API URL
const RootURL = "https://api.booli.se"

// MaxLimitResponseSize is default size of API responses
const MaxLimitResponseSize = 100

// CallerID and PrivateKey are set in init()
var callerID, privateKey string

// letters are used to create random strings
const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"

func init() {
	if callerID = os.Getenv("BOOLI_CALLER_ID"); callerID == "" {
		panic("Missing Booli callerID")
	}
	if privateKey = os.Getenv("BOOLI_PRIVATE_KEY"); privateKey == "" {
		panic("Missing Booli privateKey")
	}
}

// srand returns random string of size
func srand(size int) string {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte(letters[rand.Intn(len(letters))])
	}

	return string(buf)
}

// Query is Booli query parameters
type Query map[string]string

// Position see https://www.booli.se/p/api/referens/
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Address see https://www.booli.se/p/api/referens/
type Address struct {
	StreetAddress string `json:"streetAddress"`
}

// Region see https://www.booli.se/p/api/referens/
type Region struct {
	MunicipalityName string `json:"municipalityName"`
	CountyName       string `json:"countyName"`
}

// Location see https://www.booli.se/p/api/referens/
type Location struct {
	Position   Position `json:"position"`
	NamedAreas []string `json:"namedAreas"`
	Address    Address  `json:"address"`
	Region     Region   `json:"region"`
}

// Source see https://www.booli.se/p/api/referens/
type Source struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
}

// PublishedDate is on the "2006-01-02 15:04:05" format
type PublishedDate time.Time

// UnmarshalJSON parses Booli published date
func (j *PublishedDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = PublishedDate(t)
	return nil
}

// SoldDate is on the 2006-01-02 format
type SoldDate time.Time

// UnmarshalJSON parses Booli sold date
func (j *SoldDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = SoldDate(t)
	return nil
}

// Property see https://www.booli.se/p/api/referens/
type Property struct {
	Location            Location      `json:"location"`
	ListPrice           int           `json:"listPrice"`
	FirstPrice          int           `json:"firstPrice"`
	SoldPrice           int           `json:"soldPrice"`
	SoldDate            SoldDate      `json:"soldDate"`
	ListPriceChangeDate PublishedDate `json:"listPriceChangeDate"`
	BooliID             int           `json:"booliId"`
	Published           PublishedDate `json:"published"`
	URL                 string        `json:"url"`
	ObjectType          string        `json:"objectType"`
	Rooms               float32       `json:"rooms"`
	LivingArea          float32       `json:"livingArea"`
	PlotArea            float32       `json:"plotArea"`
	AdditionalArea      float32       `json:"additionalArea"`
	Rent                int           `json:"rent"`
	Floor               int           `json:"floor"`
	ConstructionYear    int           `json:"constructionYear"`
	Source              Source        `json:"source"`
	IsNewConstruction   int           `json:"isNewConstruction"`
	HasPatio            int           `json:"hasPatio"`
	HasBalcony          int           `json:"hasBalcony"`
	HasSolarPanels      int           `json:"hasSolarPanels"`
	HasFirePlace        int           `json:"hasFirePlace"`
	BiddingOpen         int           `json:"biddingOpen"`
	MortgageDeed        int           `json:"mortageDeed"`
	BuildingHasElevator int           `json:"buildingHasElevator"`
}

func (p Property) String() string {
	return fmt.Sprintf("%s, %s, %s", p.Location.Address.StreetAddress, p.Location.Region.MunicipalityName, p.Location.Region.CountyName)
}

// Area is an area
type Area struct {
	BooliID       int      `json:"booliId"`
	Name          string   `json:"string"`
	Types         []string `json:"types"`
	ParentBooliID int      `json:"parentBooliId"`
	ParentName    string   `json:"parentName"`
	ParentTypes   []string `json:"parentTypes"`
	FullName      string   `json:"fullName"`
}

// Pagination defines https://www.booli.se/api/#pagination
type Pagination struct {
	Count      int `json:"count"`
	TotalCount int `json:"totalCount"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
}

// Response is a generic API Response container
type Response struct {
	Pagination
	Sold     []Property `json:"sold"`     //Either this
	Areas    []Area     `json:"areas"`    // or this
	Listings []Property `json:"listings"` // or this
}

// get gets one page of response data from path.
func get(path string, params Query) ([]byte, error) {

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
	u := RootURL + "/" + path + "?" + q.Encode()
	//log.Println(u)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)

}

// paginator returns an iterator for all page responses
func paginator(resource string, params Query) chan Response {

	ch := make(chan Response)

	go func() {

		var offset int
		limit := MaxLimitResponseSize
		for {

			//Set pagination parameters
			params["limit"] = strconv.Itoa(limit)
			params["offset"] = strconv.Itoa(offset)

			data, err := get(resource, params)
			if err != nil {
				break
			}

			var resp Response
			err = json.Unmarshal(data, &resp)
			if err != nil {
				log.Printf("Could not unmarshal. %s", err)
				break
			}

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

// Sold iterates sold properties
func Sold(params Query) chan Property {
	ch := make(chan Property)
	go func() {
		for resp := range paginator(SoldResource, params) {
			for _, a := range resp.Sold {
				ch <- a
			}
		}
		close(ch)
	}()
	return ch
}

// Listings iterates listed properties
func Listings(params Query) chan Property {
	ch := make(chan Property)
	go func() {
		for resp := range paginator(ListingsResource, params) {
			for _, a := range resp.Listings {
				ch <- a
			}
		}
		close(ch)
	}()
	return ch
}

// Areas iterates areas
func Areas(params Query) chan Area {
	ch := make(chan Area)
	go func() {
		for resp := range paginator(AreasResource, params) {
			for _, a := range resp.Areas {
				ch <- a
			}
		}
		close(ch)
	}()
	return ch
}

// ImageURL returns the image URL for a property
func (p Property) ImageURL(w, h int) string {
	return fmt.Sprintf("https://bcdn.se/cache/primary_%v_%dx%d.jpg", p.BooliID, w, h)
}
