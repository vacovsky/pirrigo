package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//Settings Describes general-purpose application settings for PirriGo as derived from the configuration file (used created/populated).
type Settings struct {
	Pirri struct {
		MonitorInterval int    `json:"monitor_interval"`
		UtcOffset       int    `json:"utc_offset"`
		RainSkip        bool   `json:"rain_skip"`
		WelcomeMessage  string `json:"welcome_message"`
		UseRabbitMQ     bool   `json:"use_rabbit"`
	} `json:"pirri"`
	RabbitMQ struct {
		Server    string `json:"server"`
		User      string `json:"user"`
		Secret    string `json:"secret"`
		TaskQueue string `json:"task_queue"`
		Port      string `json:"port"`
	} `json:"rabbitmq"`
	SQL struct {
		DBType string `json:"dbtype"`
		Server string `json:"server"`
		User   string `json:"user"`
		Secret string `json:"secret"`
		Port   string `json:"port"`
		DB     string `json:"db"`
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
		Pirri        bool `json:"pirri"`
		GORM         bool `json:"gorm"`
		Settings     bool `json:"settings"`
		SimulateGPIO bool `json:"simulate_gpio"`
	} `json:"debug"`
}

func (s *Settings) parseSettingsFile() {
	confFile := "../../init/appconfig.json"
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
	setRabbitMQConnectionString()

	// set sql connection string
	setSQLConnectionString()

	// load new relic key from file
	loadNewRelicKey()

}

func setRabbitMQConnectionString() {
	if (SETTINGS.RabbitMQ.User != "" || SETTINGS.RabbitMQ.Secret != "") || SETTINGS.RabbitMQ.Server == "" {
		RMQCONNSTRING = "amqp://" + SETTINGS.RabbitMQ.User + ":" + SETTINGS.RabbitMQ.Secret + "@" + SETTINGS.RabbitMQ.Server + ":" + SETTINGS.RabbitMQ.Port + "/"
	} else {
		RMQCONNSTRING = "amqp://localhost:" + SETTINGS.RabbitMQ.Port + "/"
	}
	c := RMQCONNSTRING
	fmt.Println("RabbitMQ Connection String:", c)
}

func setSQLConnectionString() {
	SQLConnString = SETTINGS.SQL.User + ":" + SETTINGS.SQL.Secret + "@tcp(" + SETTINGS.SQL.Server + ":" + SETTINGS.SQL.Port + ")/" + SETTINGS.SQL.DB + "?parseTime=true"
	c := SQLConnString
	fmt.Println("SQL Connection String:", c)
}

func loadNewRelicKey() {
	if SETTINGS.NewRelic.Active {
		file, err := os.Open(SETTINGS.NewRelic.NewRelicLicensePath)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		key := ""
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			key = scanner.Text()
		}
		SETTINGS.NewRelic.Key = key
	}
}
