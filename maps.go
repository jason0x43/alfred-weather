package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const mapAPI = "https://maps.googleapis.com/maps/api"

type geoPos struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type geoBounds struct {
	Northeast geoPos `json:"northeast"`
	Southwest geoPos `json:"southwest"`
}

type geoLocation struct {
	AddressComponents []struct {
		LongName  string   `json:"long_name"`
		ShortName string   `json:"short_name"`
		Types     []string `json:"types"`
	} `json:"address_components"`
	FormattedAddress string `json:"formatted_address"`
	Geometry         struct {
		Bounds       geoBounds `json:"bounds"`
		Location     geoPos    `json:"location"`
		LocationType string    `json:"location_type"`
		Viewport     geoBounds `json:"viewport"`
	} `json:"geometry"`
	PlaceID string   `json:"place_id"`
	Types   []string `json:"types"`
}

type geoResults struct {
	Status  string        `json:"status"`
	Results []geoLocation `json:"results"`
}

type tzResults struct {
	DstOffset    int64  `json:"dstOffset"`
	RawOffset    int64  `json:"rawOffset"`
	Status       string `json:"status"`
	TimeZoneID   string `json:"timeZoneId"`
	TimeZoneName string `json:"timeZoneName"`
}

// Geocode is a geographic location
type Geocode struct {
	Name      string
	Latitude  float64
	Longitude float64
}

// Timezone is a timezone
type Timezone struct {
	Name      string
	ID        string
	DSTOffset int64
	UTCOffset int64
}

// Locate returns the geocode for a location
func Locate(location string) (l Geocode, err error) {
	dlog.Printf("Locating %s", location)

	url := fmt.Sprintf("%s/geocode/json", mapAPI)

	params := map[string]string{"address": location, "sensor": "false"}
	var content []byte
	if content, err = get(url, params); err != nil {
		return
	}

	var r geoResults
	if err = json.Unmarshal(content, &r); err != nil {
		return
	}

	if r.Status != "OK" {
		return l, fmt.Errorf(r.Status)
	}

	loc := &(r.Results[0])
	l.Name = loc.FormattedAddress
	l.Latitude = loc.Geometry.Location.Lat
	l.Longitude = loc.Geometry.Location.Lng

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

// GetTimezone returns the timezone for a geocode
func GetTimezone(lat float64, lng float64) (tz Timezone, err error) {
	url := fmt.Sprintf("%s/timezone/json", mapAPI)

	params := map[string]string{
		"location":  fmt.Sprintf("%f,%f", lat, lng),
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}
	var content []byte
	if content, err = get(url, params); err != nil {
		return
	}

	var t tzResults
	if err = json.Unmarshal(content, &t); err != nil {
		return
	}

	if t.Status != "OK" {
		return tz, fmt.Errorf(t.Status)
	}

	tz.Name = t.TimeZoneName
	tz.ID = t.TimeZoneID
	tz.DSTOffset = t.DstOffset
	tz.UTCOffset = t.RawOffset

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
