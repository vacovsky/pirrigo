package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func configInit() {
	loadConfig()
}

func loadConfig() {
	config := make(map[string]string)
	confFile := "app.conf"
	if len(os.Args) > 1 {
		confFile = os.Args[1]
	}

	file, err := os.Open(confFile)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c := strings.Split(scanner.Text(), "=")
		if len(c) > 1 {
			config[c[0]] = c[1]
		}
	}
	for k, v := range config {
		fmt.Println(k, "=", v)
	}
	parseConfig(config)
}

func parseConfig(config map[string]string) {
	SETTINGS = Settings{}

	if port, ok := config["port"]; ok {
		SETTINGS.HttpPort, err = strconv.Atoi(port)
	}
	if sqldbtype, ok := config["sqldbtype"]; ok {
		SETTINGS.SqlDbType = sqldbtype
	}
	if redishost, ok := config["redishost"]; ok {
		SETTINGS.RedisServer = redishost
	}
	if redisport, ok := config["redisport"]; ok {
		SETTINGS.RedisPort, err = strconv.Atoi(redisport)
	}
	if sqlserver, ok := config["sqlserver"]; ok {
		SETTINGS.SqlServer = sqlserver
	}
	if sqluser, ok := config["sqluser"]; ok {
		SETTINGS.SqlUser = sqluser
	}
	if sqlpass, ok := config["sqlpass"]; ok {
		SETTINGS.SqlPass = sqlpass
	}
	if sqldb, ok := config["sqldb"]; ok {
		SETTINGS.SqlDbName = sqldb
	}
	if rabbitserver, ok := config["rabbitserver"]; ok {
		SETTINGS.RabbitServer = rabbitserver
	}
	if rabbitport, ok := config["rabbitport"]; ok {
		SETTINGS.RabbitPort, err = strconv.Atoi(rabbitport)
	}
	if rabbituser, ok := config["rabbituser"]; ok {
		SETTINGS.RabbitUser = rabbituser
	}
	if rabbitpass, ok := config["rabbitpass"]; ok {
		SETTINGS.RabbitPass = rabbitpass
	}
	if sqlport, ok := config["sqlport"]; ok {
		SETTINGS.SqlPort = sqlport
	}
	CONNSTRING = SETTINGS.SqlUser + ":" + SETTINGS.SqlPass + "@tcp(" + SETTINGS.SqlServer + ":" + SETTINGS.SqlPort + ")/" + SETTINGS.SqlDbName
	if err != nil {
		fmt.Println(err)
	}
}

func displayConfig() {
	message := "Server available at http://localhost:" + strconv.Itoa(SETTINGS.HttpPort)
	fmt.Printf("\n" + message + "\n")
}
