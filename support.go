package main

import (
	"net/http"
	"path"
	"time"
)

const (
	unitsUS     units = "US"
	unitsMetric units = "Metric"

	serviceForecastIO   = "Forecast IO"
	serviceWunderground = "Weather Underground"
)

// Alert is a weather alert (e.g., severe thunderstorm)
type Alert struct {
	Description string    `json:"description"`
	Expires     time.Time `json:"expires"`
	URI         string    `json:"uri"`
}

// Forecast represents future weather conditions
type Forecast struct {
	Date    time.Time
	Summary string
	Details string
	Icon    string
	HiTemp  Temperature
	LowTemp Temperature
	Sunrise *time.Time
	Sunset  *time.Time
	Precip  int
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

// Temperature is a temperature in a specific unit system
type Temperature struct {
	Value float64
	Units units
}

// TimeFormats are the available time formats
var TimeFormats = []string{
	"%Y-%m-%d %H:%M",
	"%A, %B %d, %Y %I:%M%p",
	"%a, %d %b %Y %H:%M",
	"%I:%M%p on %m/%d/%Y",
	"%d.%m.%Y %H:%M",
	"%d/%m/%Y %H:%M",
}

type units string

// Weather is weather information
type Weather struct {
	Current struct {
		Summary  string
		Icon     string
		Humidity float64
		Temp     Temperature
	}
	Daily  []Forecast
	Hourly []Forecast
	Info   struct {
		Time    time.Time
		Sunrise time.Time
		Sunset  time.Time
		HiTemp  Temperature
		LowTemp Temperature
	}
	Alerts []Alert
}

var client = &http.Client{}

func getIconFile(name string) string {
	return path.Join("icons", config.Icons, name+".png")
}
