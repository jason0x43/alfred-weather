package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)
var tmIconNames = map[int64]string{
	  0: "clear",
      1000: "clear",
      1100: "mostlysunny",
      1101: "partlycloudy",
      1102: "mostlycloudy",
      1001: "cloudy",
      1103: "mostlysunny",
      2100: "fog",
      2101: "mostlysunny",
      2102: "fog",
      2103: "mostlycloudy",
      2106: "fog",
      2107: "fog",
      2108: "mostlycloudy",
      2000: "fog",
      4204: "chancerain",
      4203: "chancerain",
      4205: "chancerain",
      4000: "chancerain",
      4200: "chancerain",
      4213: "chancerain",
      4214: "chancerain",
      4215: "chanraincerain",
      4209: "rain",
      4208: "rain",
      4210: "rain",
      4001: "rain",
      4211: "rain",
      4202: "rain",
      4212: "rain",
      4201: "rain",
      5115: "flurries",
      5116: "flurries",
      5117: "flurries",
      5001: "flurries",
      5100: "flurries",
      5102: "flurries",
      5103: "flurries",
      5104: "flurries",
      5122: "flurries",
      5105: "chancesnow",
      5106: "chancesnow",
      5107: "chancesnow",
      5000: "snow",
      5101: "snow",
      5119: "snow",
      5120: "snow",
      5121: "snow",
      5110: "snow",
      5108: "snow",
      5114: "snow",
      5112: "sleet",
      6000: "sleet",
      6003: "sleet",
      6002: "sleet",
      6004: "sleet",
      6204: "sleet",
      6206: "sleet",
      6205: "sleet",
      6203: "sleet",
      6209: "sleet",
      6200: "sleet",
      6213: "sleet",
      6214: "sleet",
      6215: "sleet",
      6001: "sleet",
      6212: "sleet",
      6220: "sleet",
      6222: "sleet",
      6207: "sleet",
      6202: "sleet",
      6208: "sleet",
      6201: "sleet",
      7110: "chancesleet",
      7111: "chancesleet",
      7112: "chancesleet",
      7102: "chancesleet",
      7108: "sleet",
      7107: "sleet",
      7109: "sleet",
      7000: "sleet",
      7105: "sleet",
      7106: "sleet",
      7115: "sleet",
      7117: "sleet",
      7103: "sleet",
      7113: "sleet",
      7114: "sleet",
      7116: "sleet",
      7101: "sleet",
      8001: "tstorms",
      8003: "tstorms",
      8002: "tstorms",
      8000: "tstorms",
    
}

var tmDescriptions = map[int64]string{
	0: "Unknown",
      1000: "Clear, Sunny",
      1100: "Mostly Clear",
      1101: "Partly Cloudy",
      1102: "Mostly Cloudy",
      1001: "Cloudy",
      1103: "Partly Cloudy and Mostly Clear",
      2100: "Light Fog",
      2101: "Mostly Clear and Light Fog",
      2102: "Partly Cloudy and Light Fog",
      2103: "Mostly Cloudy and Light Fog",
      2106: "Mostly Clear and Fog",
      2107: "Partly Cloudy and Fog",
      2108: "Mostly Cloudy and Fog",
      2000: "Fog",
      4204: "Partly Cloudy and Drizzle",
      4203: "Mostly Clear and Drizzle",
      4205: "Mostly Cloudy and Drizzle",
      4000: "Drizzle",
      4200: "Light Rain",
      4213: "Mostly Clear and Light Rain",
      4214: "Partly Cloudy and Light Rain",
      4215: "Mostly Cloudy and Light Rain",
      4209: "Mostly Clear and Rain",
      4208: "Partly Cloudy and Rain",
      4210: "Mostly Cloudy and Rain",
      4001: "Rain",
      4211: "Mostly Clear and Heavy Rain",
      4202: "Partly Cloudy and Heavy Rain",
      4212: "Mostly Cloudy and Heavy Rain",
      4201: "Heavy Rain",
      5115: "Mostly Clear and Flurries",
      5116: "Partly Cloudy and Flurries",
      5117: "Mostly Cloudy and Flurries",
      5001: "Flurries",
      5100: "Light Snow",
      5102: "Mostly Clear and Light Snow",
      5103: "Partly Cloudy and Light Snow",
      5104: "Mostly Cloudy and Light Snow",
      5122: "Drizzle and Light Snow",
      5105: "Mostly Clear and Snow",
      5106: "Partly Cloudy and Snow",
      5107: "Mostly Cloudy and Snow",
      5000: "Snow",
      5101: "Heavy Snow",
      5119: "Mostly Clear and Heavy Snow",
      5120: "Partly Cloudy and Heavy Snow",
      5121: "Mostly Cloudy and Heavy Snow",
      5110: "Drizzle and Snow",
      5108: "Rain and Snow",
      5114: "Snow and Freezing Rain",
      5112: "Snow and Ice Pellets",
      6000: "Freezing Drizzle",
      6003: "Mostly Clear and Freezing drizzle",
      6002: "Partly Cloudy and Freezing drizzle",
      6004: "Mostly Cloudy and Freezing drizzle",
      6204: "Drizzle and Freezing Drizzle",
      6206: "Light Rain and Freezing Drizzle",
      6205: "Mostly Clear and Light Freezing Rain",
      6203: "Partly Cloudy and Light Freezing Rain",
      6209: "Mostly Cloudy and Light Freezing Rain",
      6200: "Light Freezing Rain",
      6213: "Mostly Clear and Freezing Rain",
      6214: "Partly Cloudy and Freezing Rain",
      6215: "Mostly Cloudy and Freezing Rain",
      6001: "Freezing Rain",
      6212: "Drizzle and Freezing Rain",
      6220: "Light Rain and Freezing Rain",
      6222: "Rain and Freezing Rain",
      6207: "Mostly Clear and Heavy Freezing Rain",
      6202: "Partly Cloudy and Heavy Freezing Rain",
      6208: "Mostly Cloudy and Heavy Freezing Rain",
      6201: "Heavy Freezing Rain",
      7110: "Mostly Clear and Light Ice Pellets",
      7111: "Partly Cloudy and Light Ice Pellets",
      7112: "Mostly Cloudy and Light Ice Pellets",
      7102: "Light Ice Pellets",
      7108: "Mostly Clear and Ice Pellets",
      7107: "Partly Cloudy and Ice Pellets",
      7109: "Mostly Cloudy and Ice Pellets",
      7000: "Ice Pellets",
      7105: "Drizzle and Ice Pellets",
      7106: "Freezing Rain and Ice Pellets",
      7115: "Light Rain and Ice Pellets",
      7117: "Rain and Ice Pellets",
      7103: "Freezing Rain and Heavy Ice Pellets",
      7113: "Mostly Clear and Heavy Ice Pellets",
      7114: "Partly Cloudy and Heavy Ice Pellets",
      7116: "Mostly Cloudy and Heavy Ice Pellets",
      7101: "Heavy Ice Pellets",
      8001: "Mostly Clear and Thunderstorm",
      8003: "Partly Cloudy and Thunderstorm",
      8002: "Mostly Cloudy and Thunderstorm",
      8000: "Thunderstorm",
    
}

const tmAPI = "https://api.Tomorrow.io/v4/timelines"

const tmUnits = "imperial"

// TomorrowIO is a weather service handle
type TomorrowIO struct {
	apiKey string
}

type tmWeather struct {
	Data struct {
		Timelines []struct {
			EndTime   string `json:"endTime"`
			Intervals []struct {
				StartTime string `json:"startTime"`
				Values    struct {
					Humidity                 float64   `json:"humidity"`
					PrecipitationProbability int64   `json:"precipitationProbability"`
					SunriseTime              string  `json:"sunriseTime"`
					SunsetTime               string  `json:"sunsetTime"`
					Temperature              float64 `json:"temperature"`
					TemperatureApparent      float64 `json:"temperatureApparent"`
					WeatherCode              int64   `json:"weatherCode"`
				} `json:"values"`
			} `json:"intervals"`
			StartTime string `json:"startTime"`
			Timestep  string `json:"timestep"`
		} `json:"timelines"`
	} `json:"data"`
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

	weather.Current.Summary = tmDescriptions[hourly.Data.Timelines[0].Intervals[0].Values.WeatherCode]
	weather.Current.Icon = tmIconNames[hourly.Data.Timelines[0].Intervals[0].Values.WeatherCode]
	weather.Current.Humidity = hourly.Data.Timelines[0].Intervals[0].Values.Humidity
	weather.Current.Temp = temperature(hourly.Data.Timelines[0].Intervals[0].Values.Temperature)
	weather.Current.ApparentTemp = temperature(hourly.Data.Timelines[0].Intervals[0].Values.TemperatureApparent)

	for j, d := range daily.Data.Timelines[0].Intervals {
		TempArr := []float64{}
		lowTempTemp := float64(0)
		highTempTemp := float64(0)
		highTemp := int(d.Values.Temperature)
		lowTemp := int(d.Values.Temperature)
		if(j==0){
			for i := 0; i <= 23; i++ {
				TempArr = append(TempArr, hourly.Data.Timelines[0].Intervals[i].Values.Temperature)
			} 
			highTempTemp, lowTempTemp = findMinAndMax(TempArr)
			highTemp = int(highTempTemp)
			lowTemp = int(lowTempTemp)
		}else if(j==1){
			for i := 23; i <= 47; i++ {
				TempArr = append(TempArr, hourly.Data.Timelines[0].Intervals[i].Values.Temperature)
			} 
			highTempTemp, lowTempTemp = findMinAndMax(TempArr)
			highTemp = int(highTempTemp)
			lowTemp = int(lowTempTemp)
		}else if(j==2){
			for i := 48; i <= 71; i++ {
				TempArr = append(TempArr, hourly.Data.Timelines[0].Intervals[i].Values.Temperature)
			} 
			highTempTemp, lowTempTemp = findMinAndMax(TempArr)
			highTemp = int(highTempTemp)
			lowTemp = int(lowTempTemp)
		}else if(j==3){
			for i := 72; i <= 95; i++ {
				TempArr = append(TempArr, hourly.Data.Timelines[0].Intervals[i].Values.Temperature)
			} 
			highTempTemp, lowTempTemp = findMinAndMax(TempArr)
			highTemp = int(highTempTemp)
			lowTemp = int(lowTempTemp)
		}else{
			
		}
		f := dailyForecast{
			Date:     parseDate(d.StartTime),
			Icon:     tmIconNames[d.Values.WeatherCode],
			Summary:  tmDescriptions[d.Values.WeatherCode],
			HighTemp: temperature(highTemp),
			LowTemp:  temperature(lowTemp),
			Sunrise:  parseTime(d.Values.SunriseTime),
			Sunset:   parseTime(d.Values.SunsetTime),
			Precip:   int(d.Values.PrecipitationProbability),
		}
		weather.Daily = append(weather.Daily, f)

		dlog.Printf("initialized precip to %d\n", f.Precip)
	}

	for _, d := range hourly.Data.Timelines[0].Intervals {
		f := hourlyForecast{
			Time:         parseTime(d.StartTime),
			Icon:         tmIconNames[d.Values.WeatherCode],
			Summary:      tmDescriptions[d.Values.WeatherCode],
			Temp:         temperature(d.Values.Temperature),
			ApparentTemp: temperature(d.Values.TemperatureApparent),
			Precip:   int(d.Values.PrecipitationProbability),
		}
		weather.Hourly = append(weather.Hourly, f)
	}

	return
}

func (f *TomorrowIO) DailyForecast(l Location) (data tmWeather, err error) {
	dlog.Printf("getting daily forecast for %#v", l)

	query := url.Values{}
	query.Set("location",fmt.Sprintf("%f", l.Latitude)+ "," + fmt.Sprintf("%f", l.Longitude))
	query.Set("fields", "precipitationProbability,temperature,temperatureApparent,weatherCode,humidity,sunriseTime,sunsetTime",)
	query.Set("units", tmUnits)
	query.Set("timesteps", "1d")
	query.Set("timezone", "America/New_York")
	query.Set("apikey", f.apiKey)

	url := fmt.Sprintf("%s?%s", tmAPI, query.Encode())
	
	dlog.Printf("getting URL Daily %s", url)

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

func (f *TomorrowIO) HourlyForecast(l Location) (data tmWeather, err error) {
	dlog.Printf("getting hourly forecast for %#v", l)

	query := url.Values{}
	query.Set("location",fmt.Sprintf("%f", l.Latitude)+ "," + fmt.Sprintf("%f", l.Longitude))
	query.Set("fields", "precipitationProbability,temperature,temperatureApparent,weatherCode,humidity",)
	query.Set("units", tmUnits)
	query.Set("timesteps", "1h")
	query.Set("timezone", "America/New_York")
	query.Set("apikey", f.apiKey)

	url := fmt.Sprintf("%s?%s", tmAPI, query.Encode())

	dlog.Printf("getting URL Hourly %s", url)
	
	var request *http.Request
	if request, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	dlog.Printf("got request 1")
	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}
	defer resp.Body.Close()
	dlog.Printf("got request 2")

	if resp.StatusCode >= 400 {
		return data, fmt.Errorf(resp.Status)
	}
	dlog.Printf("got request 3")

	err = json.NewDecoder(resp.Body).Decode(&data)
	dlog.Printf("got request 4")

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

func findMinAndMax(a []float64) (min float64, max float64) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}