package pirri

import (
	"encoding/json"
	"sync"
	"time"
)

//ManualStationRun describes the data required to trigger a "manual" Station activation from the web front end.
type ManualStationRun struct {
	StationID int
	Duration  int
}

//StationSchedule describes a scheduled activation for a Station
type StationSchedule struct {
	ID        int       `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	StartDate time.Time `sql:"DEFAULT:current_timestamp" gorm:"not null"`
	EndDate   time.Time `sql:"DEFAULT:'2025-01-01 00:00:00'" gorm:"not null"`
	Sunday    bool      `sql:"DEFAULT:false" gorm:"not null"`
	Monday    bool      `sql:"DEFAULT:false" gorm:"not null"`
	Tuesday   bool      `sql:"DEFAULT:false" gorm:"not null"`
	Wednesday bool      `sql:"DEFAULT:false" gorm:"not null"`
	Thursday  bool      `sql:"DEFAULT:false" gorm:"not null"`
	Friday    bool      `sql:"DEFAULT:false" gorm:"not null"`
	Saturday  bool      `sql:"DEFAULT:false" gorm:"not null"`
	StationID int       `gorm:"not null"`
	StartTime int       `gorm:"not null"`
	Duration  int       `sql:"DEFAULT:0" gorm:"not null"`
	Repeating bool      `sql:"DEFAULT:false" gorm:"not null"`
}

//StationHistory describes an entry in the historic Station run logs
type StationHistory struct {
	ID         int       `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	StationID  int       `gorm:"not null"`
	ScheduleID int       `sql:"DEFAULT:0" gorm:"not null"`
	Duration   int       `sql:"DEFAULT:0" gorm:"not null"`
	StartTime  time.Time `sql:"DEFAULT:current_timestamp"`
}

/*DripNode Describes a drip emitter */
type DripNode struct {
	ID        int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPH       float32
	StationID int
	Count     int
}

type Station struct {
	ID    int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO  int `gorm:"not null;unique"`
	Notes string
}

func reactToGpioMessage(m []byte) {
	task := Task{}
	json.Unmarshal(m, &task)
	task.execute()
}

//RunStatus indicates progress of a task, or if the routine is idle
type RunStatus struct {
	IsIdle     bool
	IsManual   bool
	StartTime  time.Time
	Duration   int
	ScheduleID int
	StationID  int
	Cancel     bool
}

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
