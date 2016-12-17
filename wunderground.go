package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	Description string `json:"description"`
	Expires     string `json:"expires_epoch"`
	Zones       []struct {
		Zone  string `json:"ZONE"`
		State string `json:"state"`
	} `json:"ZONES"`
}

type wundTemp struct {
	Celsius    json.Number `json:"celsius"`
	Fahrenheit json.Number `json:"fahrenheit"`
}

type wundTime struct {
	Hour   json.Number `json:"hour"`
	Minute json.Number `json:"minute"`
}

func (w *wundTime) toHoursAndMinutes() (h, m int64) {
	h, _ = w.Hour.Int64()
	m, _ = w.Minute.Int64()
	return
}

func (w *wundTime) toTime() (t time.Time) {
	now := time.Now()
	today := now.Format("1/2/2016")
	h, m := w.toHoursAndMinutes()
	loc := now.Location()
	t, _ = time.ParseInLocation("1/2/2006 15:04", fmt.Sprintf("%s %d:%d", today, h, m), loc)
	return
}

type wundEpoch struct {
	Epoch json.Number `json:"epoch"`
}

func (w *wundEpoch) toTime() (t time.Time) {
	epoch, _ := w.Epoch.Int64()
	return time.Unix(epoch, 0)
}

type wundWeather struct {
	Forecast struct {
		SimpleForecast struct {
			Daily []struct {
				Date         wundEpoch `json:"date"`
				Icon         string    `json:"icon"`
				Summary      string    `json:"conditions"`
				High         wundTemp  `json:"high"`
				Low          wundTemp  `json:"low"`
				PrecipChance int64     `json:"pop"`
			} `json:"forecastday"`
		} `json:"simpleforecast"`
		TextForecast struct {
			Daily []struct {
				TextUS     string `json:"fcttext"`
				TextMetric string `json:"fcttext_metric"`
			} `json:"forecastday"`
		} `json:"txt_forecast"`
	} `json:"forecast"`
	HourlyForecast []struct {
		Time struct {
			Epoch json.Number `json:"epoch"`
		} `json:"FCTTIME"`
		Summary      string      `json:"condition"`
		Humidity     json.Number `json:"humidity"`
		Icon         string      `json:"icon"`
		PrecipChance json.Number `json:"pop"`
		Temperature  struct {
			English json.Number `json:"english"`
			Metric  json.Number `json:"metric"`
		} `json:"temp"`
	} `json:"hourly_forecast"`
	Currently struct {
		TempC             float64     `json:"temp_c"`
		TempF             float64     `json:"temp_f"`
		Icon              string      `json:"icon"`
		Humidity          string      `json:"relative_humidity"`
		Summary           string      `json:"weather"`
		FeelsLikeC        json.Number `json:"feelslike_c"`
		FeelsLikeF        json.Number `json:"feelslike_f"`
		PrecipProbability float64     `json:"precipProbability"`
		LocalEpoch        json.Number `json:"local_epoc"`
		LocalTime         string      `json:"local_time_rfc822"`
		DisplayLocation   struct {
			City  string `json:"city"`
			State string `json:"state"`
			Full  string `json:"full"`
			ZIP   string `json:"zip"`
		} `json:"display_location"`
	} `json:"current_observation"`
	Alerts    []wundAlert `json:"alerts"`
	QueryZone string      `json:"query_zone"`
	Astronomy []struct {
		Sunrise struct {
			Date wundEpoch `json:"date"`
		} `json:"sunrise"`
		Sunset struct {
			Date wundEpoch `json:"date"`
		} `json:"sunset"`
	} `json:"astronomy10day"`
}

// NewWeatherUnderground returns a new WeatherUnderground handle
func NewWeatherUnderground(apiKey string) WeatherUnderground {
	return WeatherUnderground{apiKey: apiKey}
}

// Forecast returns the forecast for a given location
func (f *WeatherUnderground) Forecast(l Location) (weather Weather, err error) {
	dlog.Printf("getting forecast for %#v", l)

	url := fmt.Sprintf("%s/%s/conditions/alerts/hourly/astronomy10day/forecast10day/q/%f,%f.json",
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
		expires, _ := strconv.ParseInt(a.Expires, 10, 64)

		alert := alert{
			Description: a.Description,
			Expires:     time.Unix(expires, 0),
			URL:         getWundAlertURL(&w, &a),
		}
		weather.Alerts = append(weather.Alerts, alert)
	}

	if humidity, err := strconv.ParseFloat(w.Currently.Humidity[:len(w.Currently.Humidity)-1], 64); err != nil {
		weather.Current.Humidity = humidity
	}

	weather.Current.Temp = Temperature{
		Units: unitsMetric,
		Value: w.Currently.TempC,
	}

	for i, d := range w.Forecast.SimpleForecast.Daily {
		epoch, _ := d.Date.Epoch.Int64()
		highTemp, _ := d.High.Celsius.Float64()
		lowTemp, _ := d.Low.Celsius.Float64()

		f := dailyForecast{
			Date:    time.Unix(epoch, 0),
			Icon:    fromWundIconName(d.Icon),
			Precip:  int(d.PrecipChance),
			Summary: d.Summary,
			HighTemp: Temperature{
				Value: highTemp,
				Units: unitsMetric,
			},
			LowTemp: Temperature{
				Value: lowTemp,
				Units: unitsMetric,
			},
			Sunrise: w.Astronomy[i].Sunrise.Date.toTime(),
			Sunset:  w.Astronomy[i].Sunset.Date.toTime(),
		}

		weather.Daily = append(weather.Daily, f)
	}

	for _, d := range w.HourlyForecast {
		epochValue, _ := d.Time.Epoch.Int64()
		precipChance, _ := d.PrecipChance.Int64()
		temp, _ := d.Temperature.Metric.Float64()
		dtime := time.Unix(epochValue, 0)

		nt := ""
		if weather.IsAtNight(dtime) {
			nt = "nt_"
		}

		f := hourlyForecast{
			Time:   dtime,
			Icon:   nt + fromDSIconName(d.Icon),
			Precip: int(precipChance),
			Temp: Temperature{
				Value: temp,
				Units: unitsMetric,
			},
			Summary: d.Summary,
		}
		weather.Hourly = append(weather.Hourly, f)
	}

	currentTime, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", w.Currently.LocalTime)
	isNight := weather.IsAtNight(currentTime)
	nt := ""
	if isNight {
		nt = "nt_"
	}

	weather.Current.Time = currentTime
	weather.Current.Summary = w.Currently.Summary
	weather.Current.Icon = fromWundIconName(w.Currently.Icon)

	if isNight && !strings.HasPrefix(weather.Current.Icon, nt) {
		weather.Current.Icon = nt + weather.Current.Icon
	}

	return
}

func fromWundIconName(name string) string {
	if n, ok := wundIconNames[name]; ok {
		return n
	}
	return name
}

func getWundAlertURL(w *wundWeather, alert *wundAlert) (url string) {
	state := w.Currently.DisplayLocation.State
	if state != "" {
		url = fmt.Sprintf("https://www.wunderground.com/US/%s/%s.html", state, w.QueryZone)
	}
	return
}
