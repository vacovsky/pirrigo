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
	CONFIG = PirriSettings{}

	if port, ok := config["port"]; ok {
		CONFIG.HttpPort, err = strconv.Atoi(port)
	}
	if sqldbtype, ok := config["sqldbtype"]; ok {
		CONFIG.SqlDbType = sqldbtype
	}
	if redishost, ok := config["redishost"]; ok {
		CONFIG.RedisServer = redishost
	}
	if redisport, ok := config["redisport"]; ok {
		CONFIG.RedisPort, err = strconv.Atoi(redisport)
	}
	if sqlserver, ok := config["sqlserver"]; ok {
		CONFIG.SqlServer = sqlserver
	}
	if sqluser, ok := config["sqluser"]; ok {
		CONFIG.SqlUser = sqluser
	}
	if sqlpass, ok := config["sqlpass"]; ok {
		CONFIG.SqlPass = sqlpass
	}
	if sqldb, ok := config["sqldb"]; ok {
		CONFIG.SqlDbName = sqldb
	}
	if rabbitserver, ok := config["rabbitserver"]; ok {
		CONFIG.RabbitServer = rabbitserver
	}
	if rabbitport, ok := config["rabbitport"]; ok {
		CONFIG.RabbitPort, err = strconv.Atoi(rabbitport)
	}
	if rabbituser, ok := config["rabbituser"]; ok {
		CONFIG.RabbitUser = rabbituser
	}
	if rabbitpass, ok := config["rabbitpass"]; ok {
		CONFIG.RabbitPass = rabbitpass
	}
	if sqlport, ok := config["sqlport"]; ok {
		CONFIG.SqlPort = sqlport
	}
	CONNSTRING = CONFIG.SqlUser + ":" + CONFIG.SqlPass + "@tcp(" + CONFIG.SqlServer + ":" + CONFIG.SqlPort + ")/" + CONFIG.SqlDbName
	if err != nil {
		fmt.Println(err)
	}
}

func displayConfig() {
	message := "Server available at http://localhost:" + strconv.Itoa(CONFIG.HttpPort)
	fmt.Printf("\n" + message + "\n")
}
