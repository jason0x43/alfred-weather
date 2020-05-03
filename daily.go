package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
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

	heading := alfred.Item{
		Title:    "Weather for " + loc.Name,
		Subtitle: alfred.Line,
	}

	if weather.URL != "" {
		heading.AddMod(alfred.ModCmd, alfred.ItemMod{
			Subtitle: "Open this forecast in a browser",
			Arg: &alfred.ItemArg{
				Keyword: "daily",
				Mode:    alfred.ModeDo,
				Data:    alfred.Stringify(&dailyCfg{ToOpen: weather.URL}),
			},
		})
	}

	items = append(items, heading)

	deg := "F"
	if config.Units == unitsMetric {
		deg = "C"
	}

	addAlertItems(&weather, &items)

	items = append(items, alfred.Item{
		Title:    "Currently: " + weather.Current.Summary,
		Subtitle: fmt.Sprintf("%d°%s (%d°%s)", weather.Current.Temp.Int64(), deg, weather.Current.ApparentTemp.Int64(), deg),
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

		parts := []string{
			fmt.Sprintf("↓ %d°%s", entry.LowTemp.Int64(), deg),
			fmt.Sprintf("↑ %d°%s", entry.HighTemp.Int64(), deg),
		}

		if entry.Precip != -1 {
			parts = append(parts, fmt.Sprintf("☂ %d%%", entry.Precip))
		}

		parts = append(
			parts,
			fmt.Sprintf("☼ %s", entry.Sunrise.Format(config.TimeFormat)),
			fmt.Sprintf("☾ %s", entry.Sunset.Format(config.TimeFormat)),
		)

		item := alfred.Item{
			Title:    date + ": " + conditions,
			Subtitle: strings.Join(parts, "    "),
			Icon:     getIconFile(icon),
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

// Do runs the command
func (c DailyCommand) Do(data string) (out string, err error) {
	var cfg dailyCfg

	if data != "" {
		if err := json.Unmarshal([]byte(data), &cfg); err != nil {
			dlog.Printf("Error unmarshaling tag data: %v", err)
		}
	}

	if cfg.ToOpen != "" {
		dlog.Printf("opening %s", cfg.ToOpen)
		err = exec.Command("open", cfg.ToOpen).Run()
	}

	return
}

func addAlertItems(weather *Weather, items *[]alfred.Item) {
	now := time.Now()

	for _, alert := range weather.Alerts {
		if alert.Expires.After(now) {
			subtitle := fmt.Sprintf("Until %s", alert.Expires.Format(config.TimeFormat))
			expireDate := alert.Expires.Format(config.DateFormat)
			if expireDate != now.Format(config.DateFormat) {
				subtitle += fmt.Sprintf(" on %s", expireDate)
			}

			item := alfred.Item{
				Title:    alert.Description,
				Subtitle: subtitle,
				Icon:     "alert.png",
			}

			if alert.URL != "" {
				item.Arg = &alfred.ItemArg{
					Keyword: "daily",
					Mode:    alfred.ModeDo,
					Data:    alfred.Stringify(dailyCfg{ToOpen: alert.URL}),
				}
			}

			*items = append(*items, item)
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

type dailyCfg struct {
	ToOpen string
}
