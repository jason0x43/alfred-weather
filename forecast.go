package main

import (
	"strings"
	"time"

	"github.com/jason0x43/go-alfred"
)

// ForecastCommand gets a weather forecast
type ForecastCommand struct{}

// About returns information about a command
func (c ForecastCommand) About() alfred.CommandDef {
	var isEnabled bool
	if config.Service == serviceForecastIO && config.ForecastIOKey != "" {
		isEnabled = true
	} else if config.Service == serviceWunderground && config.WeatherUndergroundKey != "" {
		isEnabled = true
	}

	return alfred.CommandDef{
		Keyword:     "forecast",
		Description: "Get a forecast",
		IsEnabled:   isEnabled,
	}
}

// Items returns the items for the command
func (c ForecastCommand) Items(arg, data string) (items []alfred.Item, err error) {
	dlog.Printf("Running ForecastCommand")

	var weather Weather
	var loc Location
	expired := time.Now().Sub(cache.Time).Minutes() >= 5.0 ||
		time.Now().Format("1/2/2016") != cache.Time.Format("1/2/2016") ||
		cache.Service != config.Service

	if arg == "" && !expired {
		dlog.Printf("Using cached weather")
		weather = cache.Weather
		loc = config.Location
	} else {
		if arg != "" {
			var geo Geocode
			if geo, err = Locate(arg); err != nil {
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
	}

	items = append(items, alfred.Item{
		Title:    "Weather for " + loc.Name,
		Subtitle: alfred.Line,
	})

	for _, entry := range weather.Daily {
		var date string
		now := time.Now()
		conditions := entry.Summary
		icon := entry.Icon

		if entry.Date.Format("1/2/2006") == now.Format("1/2/2006") {
			if now.After(*entry.Sunset) {
				date = "Tonight"
				icon = "nt_" + icon
				conditions = strings.Replace(conditions, " day", " night", -1)
			} else {
				date = "Today"
			}
		} else {
			date = entry.Date.Format("Monday")
		}

		item := alfred.Item{
			Title: date + ": " + conditions,
			Icon:  getIconFile(icon),
		}

		if entry.Details != "" {
			item.Subtitle = entry.Details
		}

		items = append(items, item)
	}

	return
}
