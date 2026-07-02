package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/settings"
	"go.uber.org/zap"
)

var instance *Weather
var once sync.Once

// Weather fetches current weather via Open-Meteo (free, no API key).
type Weather struct {
	lock sync.Mutex
}

// Service returns the Weather singleton.
func Service() *Weather {
	once.Do(func() {
		instance = &Weather{lock: sync.Mutex{}}
	})
	return instance
}

type CurrentWeatherResponse struct {
	Temperature    float64 `json:"temperature_2m"`
	WindSpeed      float64 `json:"wind_speed_10m"`
	WeatherCode    int     `json:"weather_code"`
	Precipitation  float64 `json:"precipitation"`
	RelativeHumidity float64 `json:"relative_humidity_2m"`
	Status         string
}

// Current returns current weather. Uses Open-Meteo (no API key needed).
func (w *Weather) Current() CurrentWeatherResponse {
	set := settings.Service()
	logger := logging.Service()

	weather := CurrentWeatherResponse{Status: "Error"}

	if set.Weather.Latitude == 0 || set.Weather.Longitude == 0 {
		logger.LogEvent("Weather: no coordinates configured",
			zap.Float64("lat", set.Weather.Latitude),
			zap.Float64("lon", set.Weather.Longitude))
		return weather
	}

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,wind_speed_10m,weather_code,precipitation,relative_humidity_2m&temperature_unit=fahrenheit&wind_speed_unit=mph",
		set.Weather.Latitude, set.Weather.Longitude)

	resp, err := http.Get(url)
	if err != nil {
		logger.LogError("Unable to obtain weather", zap.String("error", err.Error()))
		return weather
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError("Unable to read weather response", zap.String("error", err.Error()))
		return weather
	}

	// ponytail: minimal unmarshal — only pull the current block we care about
	var raw struct {
		Current struct {
			Temperature    float64 `json:"temperature_2m"`
			WindSpeed      float64 `json:"wind_speed_10m"`
			WeatherCode    int     `json:"weather_code"`
			Precipitation  float64 `json:"precipitation"`
			RelativeHumidity float64 `json:"relative_humidity_2m"`
		} `json:"current"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		logger.LogError("Unable to unmarshal weather JSON", zap.String("error", err.Error()))
		return weather
	}

	weather.Temperature = raw.Current.Temperature
	weather.WindSpeed = raw.Current.WindSpeed
	weather.WeatherCode = raw.Current.WeatherCode
	weather.Precipitation = raw.Current.Precipitation
	weather.RelativeHumidity = raw.Current.RelativeHumidity
	weather.Status = "Success"
	return weather
}
