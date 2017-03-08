package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// http://api.wunderground.com/api/{API Key}/conditions/q/{State Abbreviation}/{City}.json

// WeatherUndergroundCurrentWeatherResponse is a JSON contract for Current Weather
type WeatherUndergroundCurrentWeatherResponse struct {
	Status   string `json:"status"`
	Response struct {
		Version        string `json:"version"`
		TermsofService string `json:"termsofService"`
		Features       struct {
			Conditions int `json:"conditions"`
		} `json:"features"`
	} `json:"response"`
	CurrentObservation struct {
		Image struct {
			URL   string `json:"url"`
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"image"`
		DisplayLocation struct {
			Full           string `json:"full"`
			City           string `json:"city"`
			State          string `json:"state"`
			StateName      string `json:"state_name"`
			Country        string `json:"country"`
			CountryIso3166 string `json:"country_iso3166"`
			Zip            string `json:"zip"`
			Magic          string `json:"magic"`
			Wmo            string `json:"wmo"`
			Latitude       string `json:"latitude"`
			Longitude      string `json:"longitude"`
			Elevation      string `json:"elevation"`
		} `json:"display_location"`
		ObservationLocation struct {
			Full           string `json:"full"`
			City           string `json:"city"`
			State          string `json:"state"`
			Country        string `json:"country"`
			CountryIso3166 string `json:"country_iso3166"`
			Latitude       string `json:"latitude"`
			Longitude      string `json:"longitude"`
			Elevation      string `json:"elevation"`
		} `json:"observation_location"`
		Estimated struct {
		} `json:"estimated"`
		StationID             string  `json:"station_id"`
		ObservationTime       string  `json:"observation_time"`
		ObservationTimeRfc822 string  `json:"observation_time_rfc822"`
		ObservationEpoch      string  `json:"observation_epoch"`
		LocalTimeRfc822       string  `json:"local_time_rfc822"`
		LocalEpoch            string  `json:"local_epoch"`
		LocalTzShort          string  `json:"local_tz_short"`
		LocalTzLong           string  `json:"local_tz_long"`
		LocalTzOffset         string  `json:"local_tz_offset"`
		Weather               string  `json:"weather"`
		TemperatureString     string  `json:"temperature_string"`
		TempF                 float64 `json:"temp_f"`
		TempC                 float64 `json:"temp_c"`
		RelativeHumidity      string  `json:"relative_humidity"`
		WindString            string  `json:"wind_string"`
		WindDir               string  `json:"wind_dir"`
		WindDegrees           int     `json:"wind_degrees"`
		WindMph               float64 `json:"wind_mph"`
		WindGustMph           int     `json:"wind_gust_mph"`
		WindKph               float64 `json:"wind_kph"`
		WindGustKph           int     `json:"wind_gust_kph"`
		PressureMb            string  `json:"pressure_mb"`
		PressureIn            string  `json:"pressure_in"`
		PressureTrend         string  `json:"pressure_trend"`
		DewpointString        string  `json:"dewpoint_string"`
		DewpointF             int     `json:"dewpoint_f"`
		DewpointC             int     `json:"dewpoint_c"`
		HeatIndexString       string  `json:"heat_index_string"`
		HeatIndexF            string  `json:"heat_index_f"`
		HeatIndexC            string  `json:"heat_index_c"`
		WindchillString       string  `json:"windchill_string"`
		WindchillF            string  `json:"windchill_f"`
		WindchillC            string  `json:"windchill_c"`
		FeelslikeString       string  `json:"feelslike_string"`
		FeelslikeF            string  `json:"feelslike_f"`
		FeelslikeC            string  `json:"feelslike_c"`
		VisibilityMi          string  `json:"visibility_mi"`
		VisibilityKm          string  `json:"visibility_km"`
		Solarradiation        string  `json:"solarradiation"`
		UV                    string  `json:"UV"`
		Precip1HrString       string  `json:"precip_1hr_string"`
		Precip1HrIn           string  `json:"precip_1hr_in"`
		Precip1HrMetric       string  `json:"precip_1hr_metric"`
		PrecipTodayString     string  `json:"precip_today_string"`
		PrecipTodayIn         string  `json:"precip_today_in"`
		PrecipTodayMetric     string  `json:"precip_today_metric"`
		Icon                  string  `json:"icon"`
		IconURL               string  `json:"icon_url"`
		ForecastURL           string  `json:"forecast_url"`
		HistoryURL            string  `json:"history_url"`
		ObURL                 string  `json:"ob_url"`
		Nowcast               string  `json:"nowcast"`
	} `json:"current_observation"`
}

func getCurrentWeather() WeatherUndergroundCurrentWeatherResponse {
	weather := WeatherUndergroundCurrentWeatherResponse{}

	if SETTINGS.Weather.WundergroundKey != "" && SETTINGS.Weather.StateAbbreviation != "" && SETTINGS.Weather.City != "" {
		weatherEndpoint := fmt.Sprintf("http://api.wunderground.com/api/%s/conditions/q/%s/%s.json", SETTINGS.Weather.WundergroundKey, SETTINGS.Weather.StateAbbreviation, SETTINGS.Weather.City)

		r, err := http.Get(weatherEndpoint)
		if err != nil {
			fmt.Sprintln("Unable to obtain weather for %s, %s.", SETTINGS.Weather.StateAbbreviation, SETTINGS.Weather.City)
		}
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		weather.Status = "Success"
		err = json.Unmarshal(body, &weather)
	} else {
		weather.Status = "Unable to connect - ensure app.conf contains location info and API key for Weather Underground."
	}
	return weather
}
