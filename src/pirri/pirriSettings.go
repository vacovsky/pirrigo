package main

type Settings struct {
	IdSettings           int
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
	UtcOffset            int
	HttpPort             string
	SqlDbType            string
	SqlServer            string
	SqlUser              string
	SqlPass              string
	SqlPort              string
	SqlDbName            string
	RedisServer          string
	RedisPort            string
	RabbitPort           string
	GormDebug            bool
	ShowSettings         bool
	SimulateGpioActivity bool
	MonitorInterval      int
	PirriDebug           bool
}
