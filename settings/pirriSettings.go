package settings

import (
	"sync"
)

var instance *Settings
var once sync.Once

// Service provides a singleton of Settings
func Service() *Settings {
	once.Do(func() {
		instance = &Settings{
			lock: sync.Mutex{},
		}
		instance.init()
	})
	return instance
}

func (s *Settings) init() {}

//Settings Describes general-purpose application settings for PirriGo as derived from the configuration file (used created/populated).
type Settings struct {
	lock  sync.Mutex
	Pirri struct {
		MonitorInterval int    `json:"monitor_interval"`
		UtcOffset       int    `json:"utc_offset"`
		RainSkip        bool   `json:"rain_skip"`
		WelcomeMessage  string `json:"welcome_message"`
		UseRabbitMQ     bool   `json:"use_rabbit"`
		DateFormat      string `json:"date_format"`
		Version         string
	} `json:"pirri"`
	RabbitMQ struct {
		Server           string `json:"server"`
		User             string `json:"user"`
		Secret           string `json:"secret"`
		TaskQueue        string `json:"task_queue"`
		Port             string `json:"port"`
		ConnectionString string
	} `json:"rabbitmq"`
	SQL struct {
		DBType           string `json:"dbtype"`
		Server           string `json:"server"`
		User             string `json:"user"`
		Secret           string `json:"secret"`
		Port             string `json:"port"`
		ConnectionString string
		DB               string `json:"db"`
	} `json:"sql"`
	NewRelic struct {
		Active              bool   `json:"active"`
		NewRelicLicensePath string `json:"license_path"` // path to a text file containing nothing but the key
		Key                 string `json:"key"`
	} `json:"newrelic"`
	Redis struct {
		Server string `json:"server"`
		Port   string `json:"port"`
	} `json:"redis"`
	Web struct {
		Port   string `json:"port"`
		User   string `json:"user"`
		Secret string `json:"secret"`
	} `json:"web"`
	Weather struct {
		Units             string `json:"units"` //imperial or metric
		Station           string `json:"station"`
		WundergroundKey   string `json:"wunderground_key"`
		StateAbbreviation string `json:"state_abbreviation"`
		City              string `json:"city"`
	} `json:"weather"`
	Debug struct {
		Pirri        bool   `json:"pirri"`
		GORM         bool   `json:"gorm"`
		Settings     bool   `json:"settings"`
		SimulateGPIO bool   `json:"simulate_gpio"`
		LogPath      string `json:"log_path"`
	} `json:"debug"`
	GPIO struct {
		CommonWire int
	}
}
