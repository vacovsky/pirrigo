package main

import (
	"fmt"
)

func showVersion() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)
}

func main() {

	fmt.Printf("Launching PirriGo v%s", VERSION)

	// load settings from the configuration file
	SETTINGS.parseSettingsFile()

	// create a log file, if missing
	Logger.NewLogHelper()
	Logger.createLogFile()

	// prep ORM for usage
	gormSetup()

	// migrate DB schema and populate with seed data
	firstRunDBSetup()

	// check if we are in local debug mode, or actually doing work
	if !SETTINGS.Debug.SimulateGPIO {
		gpioClear()
	}

	// set the common wire for powering solenoids
	setCommonWire()

	// init waitgroups for concurrent processing
	WG.Add(3)

	// Start the Web application for management of schedule etc.
	go startPirriWebApp()

	// Monitor database for pre-scheduled tasks
	go startTaskMonitor()

	// Listen for tasks to execute
	if SETTINGS.Pirri.UseRabbitMQ {
		go rabbitReceive(SETTINGS.RabbitMQ.TaskQueue)
	} else {
		go listenForTasks()
	}

	// Cleanly exit after all goroutines are finished
	WG.Wait()

	fmt.Println("Waitgroup finished - exiting!")
}
