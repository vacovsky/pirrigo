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
