package main

import (
	"time"

	"github.com/jason0x43/go-alfred"
)

// Alert is a weather alert (e.g., severe thunderstorm)
type alert struct {
	Description string    `json:"description"`
	Expires     time.Time `json:"expires"`
	URI         string    `json:"uri"`
}

// DailyForecast represents future weather conditions
type dailyForecast struct {
	Date     time.Time
	Summary  string
	Icon     string
	HighTemp Temperature
	LowTemp  Temperature
	Sunrise  time.Time
	Sunset   time.Time
	Precip   int
}

// HourlyForecast represents future weather conditions
type hourlyForecast struct {
	Date    time.Time
	Summary string
	Icon    string
	Temp    Temperature
	Precip  int
}

// Int64 returns the value of the temperature in the currently configured units as an int64
func (t *Temperature) Int64() int64 {
	if t.Units == config.Units {
		return round(t.Value)
	}
	if config.Units == unitsUS {
		return round(t.Value*(9.0/5.0) + 32.0)
	}
	return round((t.Value - 32.0) * (5.0 / 9.0))
}

// Temperature is a temperature in a specific unit system
type Temperature struct {
	Value float64
	Units units
}

type units string

// Weather is weather information
type Weather struct {
	Current struct {
		Summary  string
		Icon     string
		Humidity float64
		Temp     Temperature
		Time     time.Time
	}
	Daily  []dailyForecast
	Hourly []hourlyForecast
	Alerts []alert
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
	expired := time.Now().Sub(cache.Time).Minutes() >= 5.0 ||
		time.Now().Format("1/2/2016") != cache.Time.Format("1/2/2016") ||
		cache.Service != config.Service

	if query == "" && !expired {
		dlog.Printf("Using cached weather")
		return config.Location, cache.Weather, nil
	}

	if query != "" {
		var geo Geocode
		if geo, err = Locate(query); err != nil {
			return
		}
		loc = geo.Location()
		dlog.Printf("got location")
	} else {
		loc = config.Location
		dlog.Printf("using configured location")
	}

	if loc.Name != config.Location.Name || expired {
		switch config.Service {
		case serviceForecastIO:
			service := NewForecastIO(config.ForecastIOKey)
			if weather, err = service.Forecast(loc); err != nil {
				return
			}
		case serviceWunderground:
			service := NewWeatherUnderground(config.WeatherUndergroundKey)
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
