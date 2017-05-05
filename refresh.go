package main

import (
	"time"

	"github.com/jason0x43/go-alfred"
)

// RefreshCommand forces data to refresh
type RefreshCommand struct{}

// About returns information about a command
func (c RefreshCommand) About() alfred.CommandDef {
	return alfred.CommandDef{
		Keyword:     "refresh",
		Description: "Force forecast data to be re-downloaded",
		IsEnabled:   true,
		Arg: &alfred.ItemArg{
			Keyword: "refresh",
		},
	}
}

// Items returns the items for the command
func (c RefreshCommand) Items(arg, data string) (items []alfred.Item, err error) {
	dlog.Printf("Running RefreshCommand")

	cache.Time = time.Time{}
	if err = alfred.SaveJSON(cacheFile, &cache); err == nil {
		items = append(items, alfred.Item{
			Title:    "Refreshed!",
			Subtitle: "Data will be reloaded on the next forecast",
		})
	}

	return
}
