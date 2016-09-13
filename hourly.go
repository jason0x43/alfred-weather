package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jason0x43/go-alfred"
)

// HourlyCommand gets a weather forecast
type HourlyCommand struct{}

// About returns information about a command
func (c HourlyCommand) About() alfred.CommandDef {
	var isEnabled bool
	if config.Service == serviceForecastIO && config.ForecastIOKey != "" {
		isEnabled = true
	} else if config.Service == serviceWunderground && config.WeatherUndergroundKey != "" {
		isEnabled = true
	}

	return alfred.CommandDef{
		Keyword:     "hourly",
		Description: "Get a forecast for the next few hours",
		IsEnabled:   isEnabled,
	}
}

// Items returns the items for the command
func (c HourlyCommand) Items(arg, data string) (items []alfred.Item, err error) {
	dlog.Printf("Running HourlyCommand")

	var cfg hourlyConfig
	if data != "" {
		if err := json.Unmarshal([]byte(data), &cfg); err != nil {
			dlog.Printf("Invalid hourly config")
		}
	}

	var weather Weather
	var loc Location
	if loc, weather, err = getWeather(arg); err != nil {
		return
	}

	var startTime time.Time
	if cfg.Start != nil {
		startTime = *cfg.Start
	} else if len(weather.Hourly) > 0 {
		startTime = weather.Hourly[0].Date
	}

	items = append(items, alfred.Item{
		Title:    "Weather for " + loc.Name,
		Subtitle: alfred.Line,
		Arg: &alfred.ItemArg{
			Keyword: "daily",
		},
	})

	deg := "F"
	if config.Units == unitsMetric {
		deg = "C"
	}

	for _, entry := range weather.Hourly {
		if entry.Date.Before(startTime) {
			continue
		}

		conditions := entry.Summary
		icon := entry.Icon

		item := alfred.Item{
			Title:    entry.Date.Format("Mon "+config.TimeFormat) + ": " + conditions,
			Subtitle: fmt.Sprintf("%d°%s, %d%%", entry.Temp.Int64(), deg, entry.Precip),
			Icon:     getIconFile(icon),
		}

		items = append(items, item)
	}

	return
}

type hourlyConfig struct {
	Start *time.Time
}
