package main

import (
	"fmt"
)

func showVersion() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)
}

func main() {
	configInit()
	gormSetup()
	createJunkData()
	WG.Add(4)

	// Start the Web application for management of schedule etc.
	go startPirriWebApp()
	// Monitor database for pre-scheduled tasks
	go startTaskMonitor()
	// Listen for tasks to execute
	go rabbitReceive(SETTINGS.RabbitTaskQueue)
	// Listen for stop commands
	go rabbitReceive(SETTINGS.RabbitStopQueue)
	// Cleanly exit after all goroutines are finished
	WG.Wait()
}
