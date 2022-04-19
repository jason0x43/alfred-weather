package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var tmIconNames = map[string]string{
	"clear":               "clear",
	"mostly_clear":        "mostlysunny",
	"partly_cloudy":       "partlycloudy",
	"mostly_cloudy":       "mostlycloudy",
	"cloudy":              "cloudy",
	"fog":                 "fog",
	"fog_light":           "hazy",
	"drizzle":             "rain",
	"rain_light":          "rain",
	"rain":                "rain",
	"rain_heavy":          "rain",
	"tstorm":              "tstorms",
	"flurries":            "flurries",
	"snow_light":          "snow",
	"snow":                "snow",
	"snow_heavy":          "snow",
	"ice_pellets_light":   "sleet",
	"ice_pellets":         "sleet",
	"ice_pellets_heavy":   "sleet",
	"freezing_drizzle":    "sleet",
	"freezing_rain_light": "sleet",
	"freezing_rain":       "sleet",
	"freezing_rain_heavy": "sleet",
}

var tmDescriptions = map[string]string{
	"clear":               "Clear",
	"mostly_clear":        "Mostly sunny",
	"partly_cloudy":       "Partly cloudy",
	"mostly_cloudy":       "Mostly cloudy",
	"cloudy":              "Cloudy",
	"fog":                 "Fog",
	"fog_light":           "Light fog",
	"drizzle":             "Drizzle",
	"rain_light":          "Light rain",
	"rain":                "Rain",
	"rain_heavy":          "Heavy rain",
	"tstorm":              "Thunderstorms",
	"flurries":            "Flurries",
	"snow_light":          "Light snow",
	"snow":                "Snow",
	"snow_heavy":          "Heavy snow",
	"ice_pellets_light":   "Light sleet",
	"ice_pellets":         "Sleet",
	"ice_pellets_heavy":   "Heavy sleet",
	"freezing_drizzle":    "Freezing drizzle",
	"freezing_rain_light": "Light freezing rain",
	"freezing_rain":       "Freezing rain",
	"freezing_rain_heavy": "Heavy freezing rain",
}

const tmAPI = "https://api.Tomorrow.io/v4/weather"

const tmUnits = "si"

// TomorrowIO is a weather service handle
type TomorrowIO struct {
	apiKey string
}

type tmFloatValue struct {
	Value float64 `json:"value"`
	Units string  `json:"units"`
}

type tmIntValue struct {
	Value int    `json:"value"`
	Units string `json:"units"`
}

type tmStringValue struct {
	Value string `json:"value"`
}

type tmCurrent struct {
	Temp         tmFloatValue  `json:"temp"`
	ApparentTemp tmFloatValue  `json:"feels_like"`
	Humidity     tmFloatValue  `json:"humidity"`
	WeatherCode  tmStringValue `json:"weather_code"`
	Time         tmStringValue `json:"observation_time"`
}

type tmDaily struct {
	Temp []struct {
		Time string       `json:"observation_time"`
		Min  tmFloatValue `json:"min"`
		Max  tmFloatValue `json:"max"`
	} `json:"temp"`
	PrecipProbability tmIntValue `json:"precipitation_probability"`
	ApparentTemp      []struct {
		Time string       `json:"observation_time"`
		Min  tmFloatValue `json:"min"`
		Max  tmFloatValue `json:"max"`
	} `json:"feels_like"`
	SunriseTime tmStringValue `json:"sunrise"`
	SunsetTime  tmStringValue `json:"sunset"`
	Date        tmStringValue `json:"observation_time"`
	WeatherCode tmStringValue `json:"weather_code"`
}

type tmHourly struct {
	Temp              tmFloatValue  `json:"temp"`
	ApparentTemp      tmFloatValue  `json:"feels_like"`
	PrecipProbability tmIntValue    `json:"precipitation_probability"`
	Time              tmStringValue `json:"observation_time"`
	WeatherCode       tmStringValue `json:"weather_code"`
}

// NewTomorrowIO returns a new TomorrowIO handle
func NewTomorrowIO(apiKey string) TomorrowIO {
	return TomorrowIO{apiKey: apiKey}
}

// Forecast returns the forecast for a given location
func (f *TomorrowIO) Forecast(l Location) (weather Weather, err error) {
	dlog.Printf("getting forecast for %#v", l)

	hourly, err := f.HourlyForecast(l)
	if err != nil {
		return
	}

	daily, err := f.DailyForecast(l)
	if err != nil {
		return
	}

	current, err := f.CurrentConditions(l)
	if err != nil {
		return
	}

	weather.URL = fmt.Sprintf("%s?lat=%f&lon=%f&units=%s", tmAPI, l.Latitude, l.Longitude, tmUnits)

	weather.Current.Summary = tmDescriptions[current.WeatherCode.Value]
	weather.Current.Icon = tmIconNames[current.WeatherCode.Value]
	weather.Current.Humidity = current.Humidity.Value
	weather.Current.Temp = temperature(current.Temp.Value)
	weather.Current.ApparentTemp = temperature(current.ApparentTemp.Value)

	for _, d := range daily {
		highTemp := d.Temp[0].Max.Value
		if d.Temp[0].Max.Units == "" {
			highTemp = d.Temp[1].Max.Value
		}
		lowTemp := d.Temp[0].Min.Value
		if d.Temp[0].Min.Units == "" {
			lowTemp = d.Temp[1].Min.Value
		}

		f := dailyForecast{
			Date:     parseDate(d.Date.Value),
			Icon:     tmIconNames[d.WeatherCode.Value],
			Summary:  tmDescriptions[d.WeatherCode.Value],
			HighTemp: temperature(highTemp),
			LowTemp:  temperature(lowTemp),
			Sunrise:  parseTime(d.SunriseTime.Value),
			Sunset:   parseTime(d.SunsetTime.Value),
			Precip:   d.PrecipProbability.Value,
		}
		weather.Daily = append(weather.Daily, f)

		dlog.Printf("initialized precip to %d\n", f.Precip)
	}

	for _, d := range hourly {
		f := hourlyForecast{
			Time:         parseTime(d.Time.Value),
			Icon:         tmIconNames[d.WeatherCode.Value],
			Summary:      tmDescriptions[d.WeatherCode.Value],
			Temp:         temperature(d.Temp.Value),
			ApparentTemp: temperature(d.ApparentTemp.Value),
			Precip:       d.PrecipProbability.Value,
		}
		weather.Hourly = append(weather.Hourly, f)
	}

	return
}

func (f *TomorrowIO) DailyForecast(l Location) (data []tmDaily, err error) {
	dlog.Printf("getting daily forecast for %#v", l)

	query := url.Values{}
	query.Set("lat", fmt.Sprintf("%f", l.Latitude))
	query.Set("lon", fmt.Sprintf("%f", l.Longitude))
	query.Set("apikey", f.apiKey)
	query.Set("unit_system", tmUnits)
	query.Set("start_time", "now")
	query.Set("fields", "temp,feels_like,precipitation_probability,weather_code,sunrise,sunset")

	url := fmt.Sprintf("%s/forecast/daily?%s", tmAPI, query.Encode())

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
		return data, fmt.Errorf(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

func (f *TomorrowIO) HourlyForecast(l Location) (data []tmHourly, err error) {
	dlog.Printf("getting hourly forecast for %#v", l)

	query := url.Values{}
	query.Set("lat", fmt.Sprintf("%f", l.Latitude))
	query.Set("lon", fmt.Sprintf("%f", l.Longitude))
	query.Set("apikey", f.apiKey)
	query.Set("unit_system", tmUnits)
	query.Set("fields", "temp,feels_like,precipitation_probability,weather_code")

	url := fmt.Sprintf("%s/forecast/hourly?%s", tmAPI, query.Encode())

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
		return data, fmt.Errorf(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

func (f *TomorrowIO) CurrentConditions(l Location) (data tmCurrent, err error) {
	dlog.Printf("getting current conditions for %#v", l)

	query := url.Values{}
	query.Set("lat", fmt.Sprintf("%f", l.Latitude))
	query.Set("lon", fmt.Sprintf("%f", l.Longitude))
	query.Set("apikey", f.apiKey)
	query.Set("unit_system", tmUnits)
	query.Set("fields", "temp,feels_like,weather_code,humidity")

	url := fmt.Sprintf("%s/realtime?%s", tmAPI, query.Encode())

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
		return data, fmt.Errorf(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	return
}

func parseTime(timeStr string) time.Time {
	loc := time.Now().Location()
	date, _ := time.Parse(time.RFC3339, timeStr)
	return date.In(loc)
}

func parseDate(dateStr string) time.Time {
	date, _ := time.Parse("2006-01-02", dateStr)
	return date
}
