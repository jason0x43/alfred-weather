package main

import (
	"fmt"
	"time"

	"github.com/jason0x43/go-alfred"
)

// Alert is a weather alert (e.g., severe thunderstorm)
type alert struct {
	Description string    `json:"description"`
	Expires     time.Time `json:"expires"`
	URL         string    `json:"url"`
}

// DailyForecast represents future weather conditions
type dailyForecast struct {
	Date     time.Time
	Summary  string
	Icon     string
	HighTemp temperature
	LowTemp  temperature
	Sunrise  time.Time
	Sunset   time.Time
	Precip   int
}

// HourlyForecast represents future weather conditions
type hourlyForecast struct {
	Time         time.Time
	Summary      string
	Icon         string
	Temp         temperature
	ApparentTemp temperature
	Precip       int
}

// Int64 returns the value of the temperature in the currently configured units
// as an int64. Temperatures are assumed to be in Celsius by default
func (t temperature) Int64() int64 {
	if config.Units == unitsMetric {
		return round(float64(t))
	}
	return round(float64(t)*(9.0/5.0) + 32.0)
}

// temperature is a temperature in degrees Celsius
type temperature float64

// units identifies US or Metric units
type units string

// Weather is weather information
type Weather struct {
	Current struct {
		Summary      string
		Icon         string
		Humidity     float64
		Temp         temperature
		ApparentTemp temperature
		Time         time.Time
	}
	Daily  []dailyForecast
	Hourly []hourlyForecast
	Alerts []alert
	URL    string
}

// IsAtNight indicates whether a given time is at night
func (w *Weather) IsAtNight(t time.Time) bool {
	for i := 0; i < len(w.Daily)-1; i++ {
		if t.After(w.Daily[i].Sunset) && t.Before(w.Daily[i+1].Sunrise) {
			return true
		}
	}
	return false
}

func getWeather(query string) (loc Location, weather Weather, err error) {
	if err = validateConfig(); err != nil {
		return
	}

	expired := time.Now().Sub(cache.Time).Minutes() >= 5.0 ||
		time.Now().Format("1/2/2016") != cache.Time.Format("1/2/2016") ||
		cache.Service != config.Service

	if query == "" && !expired {
		dlog.Printf("Using cached weather")
		return config.Location, cache.Weather, nil
	}

	if query != "" {
		var geos []Geocode
		if geos, err = Locate(query); err != nil {
			return
		}
		geo := geos[0]
		loc = geo.Location()
		dlog.Printf("got location")
	} else {
		loc = config.Location
		dlog.Printf("using configured location")
	}

	if loc.Name != config.Location.Name || expired {
		switch config.Service {
		case serviceDarkSky:
			service := NewDarkSky(config.DarkSkyKey)
			if weather, err = service.Forecast(loc); err != nil {
				return
			}
		case serviceOpenWeather:
			service := NewOpenWeather(config.OpenWeatherKey)
			if weather, err = service.Forecast(loc); err != nil {
				return
			}
		}

		if loc.Name == config.Location.Name {
			cache.Service = config.Service
			cache.Time = time.Now()
			cache.Weather = weather
			if err := alfred.SaveJSON(cacheFile, &cache); err != nil {
				dlog.Printf("Unable to save cache: %v", err)
			}
		}
	} else {
		weather = cache.Weather
	}

	return
}

func validateConfig() error {
	if config.Service == "" {
		return fmt.Errorf("Please choose a service")
	}

	hasKey := false
	switch config.Service {
	case serviceDarkSky:
		hasKey = config.DarkSkyKey != ""
	case serviceOpenWeather:
		hasKey = config.OpenWeatherKey != ""
	}

	if !hasKey {
		return fmt.Errorf("Please add a API key for %s", config.Service)
	}

	if config.Location.Name == "" {
		return fmt.Errorf("Please set a default location")
	}

	return nil
}
