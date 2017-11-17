package settings

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
		instance.Pirri.Version = "4.0.0"
		instance.init()
	})
	return instance
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

func (s *Settings) init() {
	confFile := "../init/appconfig.json"
	if len(os.Args) > 1 {
		confFile = os.Args[1]
	}

	file, err := os.Open(confFile)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := os.Open(confFile)
	if err != nil {
		fmt.Println("Could not open config file\n", err)
	}

	jsonParser := json.NewDecoder(fileContent)
	if err = jsonParser.Decode(&s); err != nil {
		fmt.Println("Could not load config file. Check JSON formatting.\n", err)
	}

	// set rabbitmq connection string
	s.setRabbitMQConnectionString()

	// set sql connection string
	s.setSQLConnectionString()

	// load new relic key from file
	s.loadNewRelicKey()

}

func (s *Settings) setRabbitMQConnectionString() {
	if (s.RabbitMQ.User != "" || s.RabbitMQ.Secret != "") || s.RabbitMQ.Server == "" {
		s.RabbitMQ.ConnectionString = "amqp://" + s.RabbitMQ.User + ":" + s.RabbitMQ.Secret + "@" + s.RabbitMQ.Server + ":" + s.RabbitMQ.Port + "/"
	} else {
		s.RabbitMQ.ConnectionString = "amqp://localhost:" + s.RabbitMQ.Port + "/"
	}
	// log.LogEvent("Connecting to RabbitMQ with: " + s.RabbitMQ.ConnectionString)
}

func (s *Settings) setSQLConnectionString() {
	s.SQL.ConnectionString = s.SQL.User + ":" + s.SQL.Secret + "@tcp(" + s.SQL.Server + ":" + s.SQL.Port + ")/" + s.SQL.DB + "?parseTime=true"
}

func (s *Settings) loadNewRelicKey() {
	if s.NewRelic.Active {
		file, err := os.Open(s.NewRelic.NewRelicLicensePath)
		defer file.Close()
		if err != nil {
			// log.LogError("Unable to load New Relic license key.",
			// 	zap.String("error", err.Error()))
		}
		key := ""
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			key = scanner.Text()
		}
		s.NewRelic.Key = key
	}
}
