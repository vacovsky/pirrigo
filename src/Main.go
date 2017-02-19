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
	GormSetup()
	createJunkData()
	WG.Add(3)

	// Monitor database for pre-scheduled tasks
	go taskMonitor()

	// Listen for tasks to execute
	go rabbitReceive(SETTINGS.RabbitTaskQueue)

	// Listen for stop commands
	go rabbitReceive(SETTINGS.RabbitStopQueue)

	// cleanly exit after all goroutines are finished
	WG.Wait()
}
