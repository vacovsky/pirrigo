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

	// prep ORM for usage
	gormSetup()

	// Log Startup
	// logToFile("Starting PirriGO v"+VERSION, "")
	getLogger().LogEvent("PirriGo v" + VERSION + " starting up")

	// migrate DB schema and populate with seed data
	// TODO: make this nicer.  check before running anything.
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

	fmt.Println("Waitgroup finished - exiting!")
	// TODO: make a clean exit
	// func {	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// }
	// Cleanly exit after all goroutines are finished

	WG.Wait()

	getLogger().LogEvent("PirriGo v" + VERSION + " exiting.")
}
