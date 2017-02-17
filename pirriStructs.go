package main

import (
	"time"
)

type DripNode struct {
	//SELECT id, gph, sid, count FROM pirri.dripnodes;
	ID    int
	GPH   float32
	SID   int
	Count int
}

type GpioPins struct {
	//SELECT gpio, notes FROM pirri.gpio_pins;
	GPIO  int
	Notes float32
}

type StationHistory struct {
	//SELECT id, sid, schedule_id, duration, starttime FROM pirri.history;
	ID         int
	SID        float32
	ScheduleID int
	Duration   int
	StartTime  time.Time
}

type StationScheduleEntry struct {
	//SELECT * FROM pirri.schedule;
	ID        int
	StartDate time.Time
	EndDate   time.Time
	Sunday    bool
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Fridayb   bool
	Saturday  bool
	SID       int //station ID
	StartTime int
	Duration  int
	Repeating bool
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
	ID     int
	GPIO   int
	Notes  string
	Common bool
}
