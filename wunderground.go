package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var wundIconNames = map[string]string{
	"clear-day":           "clear",
	"clear-night":         "nt_clear",
	"partly-cloudy-day":   "partlycloudy",
	"partly-cloudy-night": "nt_partlycloudy",
	"wind":                "hazy",
}

const wundAPI = "http://api.wunderground.com/api"

// WeatherUnderground is a Forecast.io service handle
type WeatherUnderground struct {
	apiKey string
}

type wundConditions struct {
}

type wundAlert struct {
	Title   string `json:"title"`
	Expires int64  `json:"exipres"`
	URI     string `json:"uri"`
}

type wundTemp struct {
	Celsius    json.Number `json:"celsius"`
	Fahrenheit json.Number `json:"fahrenheit"`
}

type wundTime struct {
	Hour   json.Number `json:"hour"`
	Minute json.Number `json:"minute"`
}

func (w *wundTime) ToHoursAndMinutes() (h, m int64) {
	h, _ = w.Hour.Int64()
	m, _ = w.Minute.Int64()
	return
}

type wundWeather struct {
	Forecast struct {
		SimpleForecast struct {
			Daily []struct {
				Date struct {
					Epoch json.Number `json:"epoch"`
				} `json:"date"`
				Icon         string   `json:"icon"`
				Summary      string   `json:"conditions"`
				High         wundTemp `json:"high"`
				Low          wundTemp `json:"low"`
				PrecipChance int64    `json:"pop"`
			} `json:"forecastday"`
		} `json:"simpleforecast"`
		TextForecast struct {
			Daily []struct {
				TextUS     string `json:"fcttext"`
				TextMetric string `json:"fcttext_metric"`
			} `json:"forecastday"`
		} `json:"txt_forecast"`
	} `json:"forecast"`
	Currently struct {
		TempC             float64     `json:"temp_c"`
		TempF             float64     `json:"temp_f"`
		Icon              string      `json:"icon"`
		Humidity          string      `json:"relative_humidity"`
		Summary           string      `json:"weather"`
		FeelsLikeC        json.Number `json:"feelslike_c"`
		FeelsLikeF        json.Number `json:"feelslike_f"`
		PrecipProbability float64     `json:"precipProbability"`
		Time              json.Number `json:"local_epoc"`
	} `json:"current_observation"`
	Alerts   []wundAlert
	SunPhase struct {
		Sunrise wundTime `json:"sunrise"`
		Sunset  wundTime `json:"sunset"`
	} `json:"sun_phase"`
}

// NewWeatherUnderground returns a new WeatherUnderground handle
func NewWeatherUnderground(apiKey string) WeatherUnderground {
	return WeatherUnderground{apiKey: apiKey}
}

// Forecast returns the forecast for a given location
func (f *WeatherUnderground) Forecast(l Location) (weather Weather, err error) {
	dlog.Printf("getting forecast for %#v", l)

	url := fmt.Sprintf("%s/%s/conditions/alerts/astronomy/forecast10day/q/%f,%f.json",
		wundAPI, f.apiKey, l.Latitude, l.Longitude)
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

	var w wundWeather
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return
	}

	for _, a := range w.Alerts {
		alert := Alert{
			Description: a.Title,
			Expires:     time.Unix(a.Expires, 0),
			URI:         a.URI,
		}
		weather.Alerts = append(weather.Alerts, alert)
	}

	ts, _ := w.Currently.Time.Int64()

	weather.Info.Time = time.Unix(ts, 0)
	weather.Current.Summary = w.Currently.Summary
	weather.Current.Icon = fromWundIconName(w.Currently.Icon)

	if humidity, err := strconv.ParseFloat(w.Currently.Humidity[:len(w.Currently.Humidity)-1], 64); err != nil {
		weather.Current.Humidity = humidity
	}

	units := config.Units

	if units == string(unitsUS) {
		weather.Current.Temp = Temperature{
			Units: unitsUS,
			Value: w.Currently.TempF,
		}
	} else {
		weather.Current.Temp = Temperature{
			Units: unitsMetric,
			Value: w.Currently.TempC,
		}
	}

	for i, d := range w.Forecast.SimpleForecast.Daily {
		epoch, _ := d.Date.Epoch.Int64()

		f := Forecast{
			Date:    time.Unix(epoch, 0),
			Icon:    fromWundIconName(d.Icon),
			Precip:  int(d.PrecipChance * 100),
			Summary: d.Summary,
			// HiTemp: Temperature{
			// 	Value: d.TemperatureMax,
			// 	Units: units,
			// },
			// LowTemp: Temperature{
			// 	Value: d.TemperatureMin,
			// 	Units: units,
			// },
		}

		if units == string(unitsUS) {
			f.Details = w.Forecast.TextForecast.Daily[i].TextUS
		} else {
			f.Details = w.Forecast.TextForecast.Daily[i].TextMetric
		}

		now := time.Now()
		today := now.Format("1/2/2006")
		if f.Date.Format("1/2/2006") == today {
			loc := now.Location()
			h, m := w.SunPhase.Sunrise.ToHoursAndMinutes()
			sunrise, _ := time.ParseInLocation("1/2/2006 15:04", fmt.Sprintf("%s %d:%d", today, h, m), loc)
			f.Sunrise = &sunrise

			h, m = w.SunPhase.Sunset.ToHoursAndMinutes()
			sunset, _ := time.ParseInLocation("1/2/2006 15:04", fmt.Sprintf("%s %d:%d", today, h, m), loc)
			f.Sunset = &sunset
		}

		weather.Daily = append(weather.Daily, f)
	}

	return
}

func fromWundIconName(name string) string {
	if n, ok := wundIconNames[name]; ok {
		return n
	}
	return name
}
