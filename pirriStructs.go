package main

import (
	"time"
)

type DripNode struct {
	ID    int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPH   float32
	SID   int
	Count int
}

type GpioPins struct {
	GPIO  int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Notes float32
}

type StationHistory struct {
	ID         int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	SID        int
	ScheduleID int
	Duration   int
	StartTime  time.Time `sql:"DEFAULT:current_timestamp" gorm:"primary_key"`
}

type StationSchedule struct {
	ID        int       `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	StartDate time.Time `sql:"DEFAULT:current_timestamp"`
	EndDate   time.Time `sql:"DEFAULT:2025-01-01T00:00:00"`
	Sunday    bool      `sql:"DEFAULT:false"`
	Monday    bool      `sql:"DEFAULT:false"`
	Tuesday   bool      `sql:"DEFAULT:false"`
	Wednesday bool      `sql:"DEFAULT:false"`
	Thursday  bool      `sql:"DEFAULT:false"`
	Fridayb   bool      `sql:"DEFAULT:false"`
	Saturday  bool      `sql:"DEFAULT:false"`
	SID       int
	StartTime int
	Duration  int  `sql:"DEFAULT:0"`
	Repeating bool `sql:"DEFAULT:false"`
}

type PirriSettings struct {
	IdSettings        int
	OpenWeatherKey    string
	Zip               int
	WeatherUnits      string //imperial or metric
	RabbitServer      string
	RabbitUser        string
	RabbitPass        string
	PirriUsername     string
	PirriPassword     string
	AdjustForWeather  bool
	AdjustForForecast bool
	GpioOnState       int // 1 or 0
	GpioOffState      int // 1 or 0
	UseNewRelic       bool
	UtcOffset         int
	HttpPort          int
	SqlDbType         string
	SqlServer         string
	SqlUser           string
	SqlPass           string
	SqlPort           string
	SqlDbName         string
	RedisServer       string
	RedisPort         int
	RabbitPort        int
}

type Station struct {
	ID     int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO   int
	Notes  string
	Common bool `sql:"DEFAULT:false"`
}
