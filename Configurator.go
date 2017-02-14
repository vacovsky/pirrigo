package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//	"strconv"
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
	if port, ok := config["port"]; ok {
		PORT = port
	}
	if sqldbtype, ok := config["sqldbtype"]; ok {
		SQLDBTYPE = sqldbtype
	}
	if redishost, ok := config["redishost"]; ok {
		REDISHOST = redishost
	}
	if sqlserver, ok := config["sqlserver"]; ok {
		SQLSERVER = sqlserver
	}
	if sqluser, ok := config["sqluser"]; ok {
		SQLUSER = sqluser
	}
	if sqlpass, ok := config["sqlpass"]; ok {
		SQLPASS = sqlpass
	}
	if sqldb, ok := config["sqldb"]; ok {
		SQLDB = sqldb
	}
	if redisport, ok := config["redisport"]; ok {
		REDISPORT = redisport
	}
	if rabbitserver, ok := config["rabbitserver"]; ok {
		RABBITSERVER = rabbitserver
	}
	if rabbitport, ok := config["rabbitport"]; ok {
		RABBITPORT = rabbitport
	}
	if rabbituser, ok := config["rabbituser"]; ok {
		RABBITUSER = rabbituser
	}
	if rabbitpass, ok := config["rabbitpass"]; ok {
		RABBITPASS = rabbitpass
	}

}

func displayConfig() {
	message := "Server available at http://localhost:" + PORT
	fmt.Printf("\n" + message + "\n")
}
