package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var owIconNames = map[string]string{
	"01d": "clear",
	"01n": "nt_clear",
	"02d": "partlycloudy",
	"02n": "nt_partlycloudy",
	"03d": "cloudy",
	"03n": "nt_cloudy",
	"04d": "mostlycloudy",
	"04n": "nt_mostlycloudy",
	"09d": "rain",
	"09n": "nt_rain",
	"10d": "rain",
	"10n": "nt_rain",
	"11d": "tstorms",
	"11n": "nt_tstorms",
	"13d": "snow",
	"13n": "nt_snow",
	"50d": "hazy",
	"50n": "nt_hazy",
}

const owAPI = "https://api.openweathermap.org/data/2.5/onecall"

// OpenWeather is a weather service handle
type OpenWeather struct {
	apiKey string
}

type owWeather struct {
	Current struct {
		Temperature         float64 `json:"temp"`
		Humidity            float64 `json:"humidity"`
		ApparentTemperature float64 `json:"feels_like"`
		Time                int64   `json:"dt"`
		Weather             []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"current"`
	Daily []struct {
		Time         int64   `json:"dt"`
		PreciptProb  float64 `json:"pop"`
		Humidity     float64 `json:"humidity"`
		ApparentTemp struct {
			Day     float64 `json:"day"`
			Evening float64 `json:"eve"`
			Morning float64 `json:"morn"`
			Night   float64 `json:"night"`
		} `json:"feels_like"`
		Temp struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"temp"`
		SunsetTime  int64 `json:"sunset"`
		SunriseTime int64 `json:"sunrise"`
		Weather     []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"daily"`
	Hourly []struct {
		Time         int64   `json:"dt"`
		PreciptProb  float64 `json:"pop"`
		ApparentTemp float64 `json:"feels_like"`
		Humidity     float64 `json:"humidity"`
		Temp         float64 `json:"temp"`
		Weather      []struct {
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
	} `json:"hourly"`
	Timezone string `json:"timezone"`
}

// NewOpenWeather returns a new OpenWeather handle
func NewOpenWeather(apiKey string) OpenWeather {
	return OpenWeather{apiKey: apiKey}
}

// Forecast returns the forecast for a given location
func (f *OpenWeather) Forecast(l Location) (weather Weather, err error) {
	dlog.Printf("getting forecast for %#v", l)

	units := "metric"

	query := url.Values{}
	query.Set("lat", fmt.Sprintf("%f", l.Latitude))
	query.Set("lon", fmt.Sprintf("%f", l.Longitude))
	query.Set("appid", f.apiKey)
	query.Set("units", units)

	url := fmt.Sprintf("%s?%s", owAPI, query.Encode())

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

	var w owWeather
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return
	}

	weather.URL = fmt.Sprintf("%s?lat=%f&lon=%f&units=%s", owAPI, l.Latitude, l.Longitude, units)

	weather.Current.Summary = w.Current.Weather[0].Description
	weather.Current.Icon = fromOWIconName(w.Current.Weather[0].Icon)
	weather.Current.Humidity = w.Current.Humidity * 100
	weather.Current.Temp = temperature(w.Current.Temperature)
	weather.Current.ApparentTemp = temperature(w.Current.ApparentTemperature)

	for _, d := range w.Daily {
		f := dailyForecast{
			Date:     time.Unix(d.Time, 0),
			Icon:     fromOWIconName(d.Weather[0].Icon),
			Summary:  d.Weather[0].Description,
			HighTemp: temperature(d.Temp.Max),
			LowTemp:  temperature(d.Temp.Min),
			Sunrise:  time.Unix(d.SunriseTime, 0),
			Sunset:   time.Unix(d.SunsetTime, 0),

			// OpenWeather doesn't support precip chance
			Precip: int(d.PreciptProb * 100),
		}
		weather.Daily = append(weather.Daily, f)
	}

	for _, d := range w.Hourly {
		f := hourlyForecast{
			Time:         time.Unix(d.Time, 0),
			Icon:         fromOWIconName(d.Weather[0].Icon),
			Summary:      d.Weather[0].Description,
			Temp:         temperature(d.Temp),
			ApparentTemp: temperature(d.ApparentTemp),
			Precip: 	 int(d.PreciptProb * 100),
		}
		weather.Hourly = append(weather.Hourly, f)
	}

	return
}

func fromOWIconName(name string) string {
	if n, ok := owIconNames[name]; ok {
		return n
	}
	return name
}
