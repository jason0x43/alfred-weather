package main

import (
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	unitsUS     units = "US"
	unitsMetric units = "Metric"

	serviceDarkSky      = "Dark Sky"
	serviceWunderground = "Weather Underground"
)

func round(val float64) int64 {
	return int64(val + 0.5)
}

// Location is a named location
type Location struct {
	Latitude  float64
	Longitude float64
	Timezone  string
	ShortName string
	Name      string
}

// Service is a forecasting service
type Service interface {
	About() string
	Forecast(Location, map[string]string) (Weather, error)
}

// TimeFormats are the available time formats
var TimeFormats = []string{
	"15:04",
	"3:04pm",
}

// DateFormats are the available time formats
var DateFormats = []string{
	"2006-1-2",
	"Mon, Jan 2, 2006",
	"Mon, 2 Jan 2006",
	"1/2/2006",
	"2.1.2006",
	"2/1/2006",
}

var client = &http.Client{}

func getIconFile(name string) string {
	icon := path.Join("icons", config.Icons, name+".png")
	if _, err := os.Stat(icon); err != nil {
		if strings.HasPrefix(name, "nt_") {
			return path.Join("icons", config.Icons, name[3:]+".png")
		}
	}
	return icon
}
