package main

//Settings Describes general-purpose application settings for PirriGo as derived from the configuration file (used created/populated).
type Settings struct {
	IDSettings           int
	OpenWeatherKey       string
	Zip                  int
	WeatherUnits         string //imperial or metric
	RabbitServer         string
	RabbitUser           string
	RabbitPass           string
	RabbitTaskQueue      string
	RabbitStopQueue      string
	PirriUsername        string
	PirriPassword        string
	AdjustForWeather     bool
	AdjustForForecast    bool
	GpioOnState          int // 1 or 0
	GpioOffState         int // 1 or 0
	UseNewRelic          bool
	NewRelicLicense      string
	NewRelicLicensePath  string // path to a text file containing nothing but the key
	UtcOffset            int
	HTTPPort             string
	SQLDBType            string
	SQLServer            string
	SQLUser              string
	SQLPass              string
	SQLPort              string
	SQLDbName            string
	RedisServer          string
	RedisPort            string
	RabbitPort           string
	GormDebug            bool
	ShowSettings         bool
	SimulateGpioActivity bool
	MonitorInterval      int
	PirriDebug           bool
	WebUser              string
	WebPassword          string
	WeatherStation       string
	WundergroundKey      string
	StateAbbreviation    string
	City                 string
}
