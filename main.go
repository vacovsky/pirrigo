package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/pirri"
)

func main() {
	logger := logging.Service()
	logger.LogEvent("PirriGo starting up!")

	log.Println("PIRRIGO_WEB_PORT:", os.Getenv("PIRRIGO_WEB_PORT"))
	log.Println("PIRRIGO_DB_TYPE:", os.Getenv("PIRRIGO_DB_TYPE"))
	log.Println("PIRRIGO_DB_PATH:", os.Getenv("PIRRIGO_DB_PATH"))
	log.Println("PIRRIGO_SIMULATE_GPIO:", os.Getenv("PIRRIGO_SIMULATE_GPIO"))
	log.Println("PIRRIGO_LOG_LOCATION:", os.Getenv("PIRRIGO_LOG_LOCATION"))
	log.Println("PIRRIGO_DB_LOGMODE:", os.Getenv("PIRRIGO_DB_LOGMODE"))
	log.Println("PIRRIGO_UTC_OFFSET:", os.Getenv("PIRRIGO_UTC_OFFSET"))
	log.Println("PIRRIGO_USERNAME:", os.Getenv("PIRRIGO_USERNAME"))
	log.Println("PIRRIGO_PASSWORD:", os.Getenv("PIRRIGO_PASSWORD"))

	// initialize dependencies
	data.Service()

	// migrate DB schema and populate with seed data
	migrateDataSchema()

	// check if we are in local debug mode, or actually doing work.
	// If not debug, reset the GPIO state
	if strings.ToLower(os.Getenv("PIRRIGO_SIMULATE_GPIO")) != "true" {
		pirri.GPIOClear()
	}

	// set the common wire for powering solenoids
	pirri.SetCommonWire()

	// init waitgroups for concurrent processing
	pirri.WG.Add(3)

	// Start the Web application for management of schedule etc.
	go pirri.StartPirriWebApp()

	// Monitor database for pre-scheduled tasks
	go pirri.StartTaskMonitor()

	// Listen for tasks to execute
	go pirri.ListenForTasks()

	go listenForExit()

	pirri.WG.Wait()
	fmt.Println("Exit key received - exiting!")
}

func listenForExit() {
	log.Println("=================== PRESS <ENTER> KEY TO EXIT ===================")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	logging.Service().LogEvent("PirriGo exiting due to the exit key being pressed.  You did this...")
	pirri.WG.Done()
	os.Exit(0)
}
