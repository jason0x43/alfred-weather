package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var dsIconNames = map[string]string{
	"clear-day":           "clear",
	"clear-night":         "nt_clear",
	"partly-cloudy-day":   "partlycloudy",
	"partly-cloudy-night": "nt_partlycloudy",
	"wind":                "hazy",
}

const dsAPI = "https://api.darksky.net/forecast"

// DarkSky is a Forecast.io service handle
type DarkSky struct {
	apiKey string
}

type dsConditions struct {
	Temperature         float64 `json:"temperature"`
	Icon                string  `json:"icon"`
	Humidity            float64 `json:"humidity"`
	Summary             string  `json:"summary"`
	ApparentTemperature float64 `json:"apparentTemperature"`
	PrecipProbability   float64 `json:"precipProbability"`
	Time                int64   `json:"time"`
}

type dsWeather struct {
	Daily struct {
		Icon string `json:"icon"`
		Data []struct {
			PrecipType        string  `json:"precipType"`
			TempMin           float64 `json:"temperatureMin"`
			TempMax           float64 `json:"temperatureMax"`
			Summary           string  `json:"summary"`
			SunsetTime        int64   `json:"sunsetTime"`
			SunriseTime       int64   `json:"sunriseTime"`
			PrecipProbability float64 `json:"precipProbability"`
			Icon              string  `json:"icon"`
			Time              int64   `json:"time"`
		} `json:"data"`
		Summary string `json:"summary"`
	} `json:"daily"`
	Hourly struct {
		Icon    string `json:"icon"`
		Summary string `json:"summary"`
		Data    []struct {
			ApparentTemp      float64 `json:"apparentTemperature"`
			Humidity          float64 `json:"humidity"`
			Icon              string  `json:"icon"`
			PrecipProbability float64 `json:"precipProbability"`
			Summary           string  `json:"summary"`
			Temp              float64 `json:"temperature"`
			Time              int64   `json:"time"`
		} `json:"data"`
	} `json:"hourly"`
	Currently dsConditions `json:"currently"`
	Flags     struct {
		Units string `json:"units"`
	} `json:"flags"`
	Alerts []struct {
		Title   string `json:"title"`
		Expires int64  `json:"exipres"`
		URI     string `json:"uri"`
	} `json:"alerts"`
}

// NewDarkSky returns a new DarkSky handle
func NewDarkSky(apiKey string) DarkSky {
	return DarkSky{apiKey: apiKey}
}

// Forecast returns the forecast for a given location
func (f *DarkSky) Forecast(l Location) (weather Weather, err error) {
	dlog.Printf("getting forecast for %#v", l)

	query := url.Values{}
	query.Set("exclude", "minutely")

	if config.Units == unitsUS {
		query.Set("units", "us")
	} else {
		query.Set("units", "si")
	}

	url := fmt.Sprintf("%s/%s/%f,%f?%s", dsAPI, f.apiKey, l.Latitude, l.Longitude, query.Encode())

	dlog.Printf("getting URL %s", url)

	var request *http.Request
	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return weather, fmt.Errorf(resp.Status)
	}

	var w dsWeather
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return
	}

	for _, a := range w.Alerts {
		alert := alert{
			Description: a.Title,
			Expires:     time.Unix(a.Expires, 0),
			URL:         a.URI,
		}
		weather.Alerts = append(weather.Alerts, alert)
	}

	units := w.Flags.Units

	weather.Current.Summary = w.Currently.Summary
	weather.Current.Icon = fromDSIconName(w.Currently.Icon)
	weather.Current.Humidity = w.Currently.Humidity * 100
	weather.Current.Temp = fromDSTemp(w.Currently.Temperature, units)

	for _, d := range w.Daily.Data {
		f := dailyForecast{
			Date:     time.Unix(d.Time, 0),
			Icon:     fromDSIconName(d.Icon),
			Precip:   int(d.PrecipProbability * 100),
			Summary:  d.Summary,
			HighTemp: fromDSTemp(d.TempMax, units),
			LowTemp:  fromDSTemp(d.TempMin, units),
			Sunrise:  time.Unix(d.SunriseTime, 0),
			Sunset:   time.Unix(d.SunsetTime, 0),
		}
		weather.Daily = append(weather.Daily, f)
	}

	for _, d := range w.Hourly.Data {
		f := hourlyForecast{
			Time:    time.Unix(d.Time, 0),
			Icon:    fromDSIconName(d.Icon),
			Precip:  int(d.PrecipProbability * 100),
			Summary: d.Summary,
			Temp:    fromDSTemp(d.Temp, units),
		}
		weather.Hourly = append(weather.Hourly, f)
	}

	return
}

func fromDSIconName(name string) string {
	if n, ok := dsIconNames[name]; ok {
		return n
	}
	return name
}

func fromDSTemp(temp float64, units string) temperature {
	if units == "si" {
		return temperature(temp)
	}
	return temperature((temp - 32.0) * (5.0 / 9.0))
}
