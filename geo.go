package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const mapAPI = "https://nominatim.openstreetmap.org/search"

type geoLocation struct {
	FormattedAddress string `json:"display_name"`
	Lat              string `json:"lat"`
	Lng              string `json:"lon"`
}

type geoResults []geoLocation

// Geocode is a geographic location
type Geocode struct {
	Name      string
	Latitude  float64
	Longitude float64
}

// Locate returns the possible geocodes for a location
func Locate(location string) (l []Geocode, err error) {
	dlog.Printf("Locating %s", location)

	url := fmt.Sprintf("%s", mapAPI)

	params := map[string]string{"q": location, "format": "json", "addressDetails": "1"}
	var content []byte
	if content, err = get(url, params); err != nil {
		return
	}

	dlog.Printf("Gt results: %s", content)

	var r geoResults
	if err = json.Unmarshal(content, &r); err != nil {
		return
	}

	for _, res := range r {
		loc := &(res)
		var gc Geocode
		gc.Name = loc.FormattedAddress
		gc.Latitude, _ = strconv.ParseFloat(loc.Lat, 64)
		gc.Longitude, _ = strconv.ParseFloat(loc.Lng, 64)
		l = append(l, gc)
	}

	return
}

// Location converts a Geocode to a Location
func (g *Geocode) Location() (l Location) {
	l.Latitude = g.Latitude
	l.Longitude = g.Longitude
	l.Name = g.Name
	l.ShortName = g.Name
	return
}

func get(url string, params map[string]string) (data []byte, err error) {
	var request *http.Request
	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}

	if params != nil {
		values := request.URL.Query()
		for k, v := range params {
			values.Add(k, v)
		}
		request.URL.RawQuery = values.Encode()
	}

	dlog.Printf("request query: %s", request.URL.RawQuery)

	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		if data != nil {
			dlog.Printf("Error getting data: %s", data)
		}
		err = fmt.Errorf(resp.Status)
	}

	return
}
