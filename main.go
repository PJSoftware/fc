package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	api "github.com/pjsoftware/go-api"
)

const apiKeyVar = "API_KEY"

const (
	apiUrl = "http://api.weatherapi.com/v1"
	epForecast = "forecast.json"
	other = "days=5&aqi=no&alerts=no"
)

func main() {
	loc := "Brisbane"
	if len(os.Args) > 1 {
		loc = os.Args[1]
	}
	apiKey := readAPIKey()
	forecast := retrieveForecast(apiKey, loc)

	location := forecast.Location
	current := forecast.Current
	hours := forecast.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s, %s: %.1fC (feels like %.1fC) - %s\n", 
		location.Name, 
		location.Country, 
		current.TempC, 
		current.FeelslikeC, 
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(int64(hour.TimeEpoch), 0)
		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s - %.1fC, %d%%, %s\n",
			date.Format("20060102 15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain <= 20 {
			fmt.Print(message)
		} else if hour.ChanceOfRain <= 50 {
			color.Yellow(message)
		} else {
			color.Red(message)
		}
	}
}

func readAPIKey() string {
	env, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	if _, ok := env[apiKeyVar]; !ok {
		panic("cannot read " + apiKeyVar + " from .env file")
	}

	apiKey := env[apiKeyVar]
	if apiKey == "xxxx" {
		panic(apiKeyVar + " must be set in .env file")
	}

	return apiKey
}

func retrieveForecast(apiKey string, loc string) Forecast {
	weather := api.New(apiUrl)
	ep := weather.NewEndpoint(epForecast)
	req := ep.NewRequest()
	req.AddQuery("key",apiKey)
	req.AddQuery("q",loc)
	req.AddQueryInt("days",2)
	req.AddQuery("aqi","no")
	req.AddQuery("alerts","no")
	
	res, err := req.GET()

	if err != nil {
		panic(err)
	}

	if res.Status != 200 {
		panic(fmt.Sprintf("error communicating with Weather API: status code %d", res.Status))
	}

	var forecast Forecast
	err = json.Unmarshal(res.Body, &forecast)
	if err != nil {
		panic(err)
	}

	return forecast
}