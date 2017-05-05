package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/jason0x43/go-alfred"
)

var cacheFile string
var configFile string
var workflow alfred.Workflow

type configStruct struct {
	Service               string   `desc:"Service to use"`
	DarkSkyKey            string   `desc:"Your API key for Dark Sky"`
	WeatherUndergroundKey string   `desc:"Your API key for Weather Underground"`
	Icons                 string   `desc:"Icon set"`
	DateFormat            string   `desc:"Date format"`
	TimeFormat            string   `desc:"Time format"`
	Location              Location `desc:"Default location"`
	Units                 units    `desc:"Units"`
}

var config configStruct

var cache struct {
	Weather Weather
	Time    time.Time
	Service string
}

var dlog = log.New(os.Stderr, "[weather] ", log.LstdFlags)

func main() {
	var err error

	if workflow, err = alfred.OpenWorkflow(".", true); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	workflow.UpdateIcon = "notice.png"

	configFile = path.Join(workflow.DataDir(), "config.json")
	cacheFile = path.Join(workflow.CacheDir(), "cache.json")

	dlog.Println("Using config file", configFile)
	dlog.Println("Using cache file", cacheFile)

	if err := alfred.LoadJSON(configFile, &config); err == nil {
		dlog.Println("loaded config")
	}

	if err := alfred.LoadJSON(cacheFile, &cache); err == nil {
		dlog.Println("loaded cache")
	}

	if config.TimeFormat == "" {
		config.TimeFormat = TimeFormats[0]
	}

	if config.DateFormat == "" {
		config.DateFormat = DateFormats[0]
	}

	if config.Icons == "" {
		config.Icons = "grzanka"
	}

	if config.Units == "" {
		config.Units = unitsUS
	}

	commands := []alfred.Command{
		DailyCommand{},
		HourlyCommand{},
		OptionsCommand{},
		RefreshCommand{},
	}

	workflow.Run(commands)
}
