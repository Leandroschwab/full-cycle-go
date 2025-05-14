package services

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
)

// Make this variable to allow replacing in tests
var weatherAPIBaseURL = "http://api.weatherapi.com/v1/current.json"

type WeatherResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

func GetTemperature(city, state string) (float64, float64, float64, error) {
	if city == "" || state == "" {
		return 0, 0, 0, fmt.Errorf("invalid city or state")
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, 0, 0, fmt.Errorf("WEATHER_API_KEY environment variable not set")
	}

	fmt.Print("Fetching temperature for city: ", city, " and state: ", state, "\n")

	location := fmt.Sprintf("%s, %s", city, state)
	encodedLocation := url.QueryEscape(location)
	requestURL := fmt.Sprintf("%s?key=%s&q=%s", weatherAPIBaseURL, apiKey, encodedLocation)

	resp, err := http.Get(requestURL)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to fetch temperature: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, 0, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to decode weather response: %v", err)
	}

	tempK := math.Round(weatherResponse.Current.TempC*10+273.15) / 10
	fmt.Printf("Weather API response: %+v, Calculated Kelvin: %.1f\n", weatherResponse, tempK)
	return weatherResponse.Current.TempC, weatherResponse.Current.TempF, tempK, nil
}
