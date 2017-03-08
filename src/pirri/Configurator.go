package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// func configInit() {
// 	loadConfig()
// }

// func loadConfig() {
// 	config := make(map[string]string)
// 	confFile := "../../init/app.conf"
// 	if len(os.Args) > 1 {
// 		confFile = os.Args[1]
// 	}

// 	file, err := os.Open(confFile)
// 	defer file.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		c := strings.Split(scanner.Text(), "=")
// 		if len(c) > 1 {
// 			config[c[0]] = c[1]
// 		}
// 	}
// 	if SETTINGS.ShowSettings {
// 		for k, v := range config {
// 			fmt.Println(k, "=", v)
// 		}
// 	}
// 	parseConfig(config)
// }

// func loadNewRelicKey(path string) string {
// 	file, err := os.Open(path)
// 	defer file.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	key := ""
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		key = scanner.Text()
// 	}
// 	return key
// }

// func parseConfig(config map[string]string) {
// 	SETTINGS = Settings{}
// 	if webuser, ok := config["webuser"]; ok {
// 		SETTINGS.WebUser = webuser
// 	}
// 	if webpass, ok := config["webpass"]; ok {
// 		SETTINGS.WebPassword = webpass
// 	}
// 	if weatherstation, ok := config["weatherstation"]; ok {
// 		SETTINGS.WeatherStation = weatherstation
// 	}
// 	if utcoffset, ok := config["utcoffset"]; ok {
// 		SETTINGS.UtcOffset, ERR = strconv.Atoi(utcoffset)
// 	}
// 	if port, ok := config["port"]; ok {
// 		SETTINGS.HTTPPort = port
// 	}
// 	if redishost, ok := config["redishost"]; ok {
// 		SETTINGS.RedisServer = redishost
// 	}
// 	if redisport, ok := config["redisport"]; ok {
// 		SETTINGS.RedisPort = redisport
// 	}
// 	if sqlserver, ok := config["sqlserver"]; ok {
// 		SETTINGS.SQLServer = sqlserver
// 	}
// 	if sqluser, ok := config["sqluser"]; ok {
// 		SETTINGS.SQLUser = sqluser
// 	}
// 	if sqlpass, ok := config["sqlpass"]; ok {
// 		SETTINGS.SQLPass = sqlpass
// 	}
// 	if sqldbname, ok := config["sqldb"]; ok {
// 		SETTINGS.SQLDbName = sqldbname
// 	}
// 	if sqldbtype, ok := config["sqldbtype"]; ok {
// 		SETTINGS.SQLDBType = sqldbtype
// 	}
// 	if rabbitserver, ok := config["rabbitserver"]; ok {
// 		SETTINGS.RabbitServer = rabbitserver
// 	}
// 	if rabbitport, ok := config["rabbitport"]; ok {
// 		SETTINGS.RabbitPort = rabbitport
// 	}
// 	if rabbituser, ok := config["rabbituser"]; ok {
// 		SETTINGS.RabbitUser = rabbituser
// 	}
// 	if rabbitpass, ok := config["rabbitpass"]; ok {
// 		SETTINGS.RabbitPass = rabbitpass
// 	}
// 	if sqlport, ok := config["sqlport"]; ok {
// 		SETTINGS.SQLPort = sqlport
// 	}
// 	if gormdebug, ok := config["gormdebug"]; ok {
// 		SETTINGS.GormDebug, ERR = strconv.ParseBool(gormdebug)
// 	}
// 	if showsettings, ok := config["showsettings"]; ok {
// 		SETTINGS.ShowSettings, ERR = strconv.ParseBool(showsettings)
// 	}
// 	if simulategpioactivity, ok := config["simulategpioactivity"]; ok {
// 		SETTINGS.SimulateGpioActivity, ERR = strconv.ParseBool(simulategpioactivity)
// 	}
// 	if monitorinterval, ok := config["monitorinterval"]; ok {
// 		SETTINGS.MonitorInterval, ERR = strconv.Atoi(monitorinterval)
// 	}
// 	if taskqueue, ok := config["taskqueue"]; ok {
// 		SETTINGS.RabbitTaskQueue = taskqueue
// 	}
// 	if stopqueue, ok := config["stopqueue"]; ok {
// 		SETTINGS.RabbitStopQueue = stopqueue
// 	}
// 	if pirridebug, ok := config["pirridebug"]; ok {
// 		SETTINGS.PirriDebug, ERR = strconv.ParseBool(pirridebug)
// 	}
// 	if wundergroundkey, ok := config["wundergroundkey"]; ok {
// 		SETTINGS.WundergroundKey = wundergroundkey
// 	}
// 	if stateabbreviation, ok := config["stateabbreviation"]; ok {
// 		SETTINGS.StateAbbreviation = stateabbreviation
// 	}
// 	if city, ok := config["city"]; ok {
// 		SETTINGS.City = city
// 	}

// 	if (SETTINGS.RabbitUser != "" || SETTINGS.RabbitPass != "") || SETTINGS.RabbitServer == "" {
// 		RMQCONNSTRING = "amqp://" + SETTINGS.RabbitUser + ":" + SETTINGS.RabbitPass + "@" + SETTINGS.RabbitServer + ":" + SETTINGS.RabbitPort + "/"
// 	} else {
// 		RMQCONNSTRING = "amqp://localhost:5672/"
// 	}
// 	SQLCONNSTRING = SETTINGS.SQLUser + ":" + SETTINGS.SQLPass + "@tcp(" + SETTINGS.SQLServer + ":" + SETTINGS.SQLPort + ")/" + SETTINGS.SQLDbName + "?parseTime=true"
// 	if SETTINGS.ShowSettings {
// 		fmt.Println(SQLCONNSTRING)
// 		fmt.Println(RMQCONNSTRING)
// 	}
// 	if usenewrelic, ok := config["usenewrelic"]; ok {
// 		SETTINGS.UseNewRelic, ERR = strconv.ParseBool(usenewrelic)
// 	}
// 	if SETTINGS.UseNewRelic {
// 		if newreliclicensepath, ok := config["newreliclicensepath"]; ok {
// 			SETTINGS.NewRelicLicensePath = newreliclicensepath
// 		}
// 	}
// 	if ERR != nil {
// 		failOnError(ERR, ERR.Error())
// 	}
// }

// func displayConfig() {
// 	message := "Server available at http://localhost:" + SETTINGS.HTTPPort
// 	fmt.Printf("\n" + message + "\n")
// }
