package main

import (
	"fmt"
)

func showVersion() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)
}

func main() {
	SETTINGS.parseSettingsFile()
	// parseSettingsFile()
	gormSetup()
	firstRunDBSetup()
	WG.Add(3)

	setCommonWire()

	// Start the Web application for management of schedule etc.
	go startPirriWebApp()
	// Monitor database for pre-scheduled tasks
	go startTaskMonitor()
	// Listen for tasks to execute
	go rabbitReceive(SETTINGS.RabbitMQ.TaskQueue)
	// Cleanly exit after all goroutines are finished
	WG.Wait()

	fmt.Println("Waitgroup finished - exiting!")
}
