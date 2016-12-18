package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jason0x43/go-alfred"
)

// DailyCommand gets a weather forecast
type DailyCommand struct{}

// About returns information about a command
func (c DailyCommand) About() alfred.CommandDef {
	return alfred.CommandDef{
		Keyword:     "daily",
		Description: "Get a forecast for the next few days",
		IsEnabled:   true,
	}
}

// Items returns the items for the command
func (c DailyCommand) Items(arg, data string) (items []alfred.Item, err error) {
	dlog.Printf("Running DailyCommand")

	var weather Weather
	var loc Location
	if loc, weather, err = getWeather(arg); err != nil {
		return
	}

	items = append(items, alfred.Item{
		Title:    "Weather for " + loc.Name,
		Subtitle: alfred.Line,
	})

	deg := "F"
	if config.Units == unitsMetric {
		deg = "C"
	}

	items = append(items, getAlertItems(&weather)...)

	items = append(items, alfred.Item{
		Title:    "Currently: " + weather.Current.Summary,
		Subtitle: fmt.Sprintf("%d°%s", weather.Current.Temp.Int64(), deg),
		Icon:     getIconFile(weather.Current.Icon),
		Arg: &alfred.ItemArg{
			Keyword: "hourly",
			Data:    alfred.Stringify(&hourlyConfig{Start: &weather.Current.Time}),
		},
	})

	for _, entry := range weather.Daily {
		var date string
		now := time.Now()
		conditions := entry.Summary
		icon := entry.Icon

		if entry.Date.Format("1/2/2006") == now.Format("1/2/2006") {
			if weather.IsAtNight(now) {
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
			Subtitle: fmt.Sprintf("%d/%d°%s    ☂ %d%%    ☼ %s    ☾ %s",
				entry.HighTemp.Int64(), entry.LowTemp.Int64(), deg,
				entry.Precip, entry.Sunrise.Format(config.TimeFormat),
				entry.Sunset.Format(config.TimeFormat)),
			Icon: getIconFile(icon),
		}

		if hasHourly(weather, entry.Date) {
			item.Arg = &alfred.ItemArg{
				Keyword: "hourly",
				Data:    alfred.Stringify(&hourlyConfig{Start: &entry.Sunrise}),
			}
		}

		items = append(items, item)
	}

	return
}

func getAlertItems(weather *Weather) (items []alfred.Item) {
	now := time.Now()

	for _, alert := range weather.Alerts {
		if alert.Expires.After(now) {
			subtitle := fmt.Sprintf("Until %s", alert.Expires.Format(config.TimeFormat))
			expireDate := alert.Expires.Format(config.DateFormat)
			if expireDate != now.Format(config.DateFormat) {
				subtitle += fmt.Sprintf(" on %s", expireDate)
			}

			item := alfred.Item{
				Title:    fmt.Sprintf("Alert: %s", alert.Description),
				Subtitle: subtitle,
				Icon:     "alert.png",
			}

			if alert.URL != "" {
				item.Arg = &alfred.ItemArg{
					Keyword: "daily",
					Mode:    alfred.ModeDo,
					Data:    alfred.Stringify(optionsCfg{ToOpen: alert.URL}),
				}
			}

			items = append(items, item)
		}
	}

	return
}

func hasHourly(weather Weather, date time.Time) bool {
	target := date.Format("2006-01-02")
	for i := range weather.Hourly {
		if weather.Hourly[i].Time.Format("2006-01-02") == target {
			return true
		}
	}
	return false
}
