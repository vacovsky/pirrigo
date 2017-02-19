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
	WG.Add(3)

	//	t := Task{Station, ScheduleID: 4, Duration: 15}
	//	task, _ := json.Marshal(t)
	//	//	failOnError(ERR, ERR.Error())

	// TODO this is for testing only
	//	RabbitSend(SETTINGS.RabbitTaskQueue, string(task))
	// end test

	// Check schedule table every minute for tasks
	go TaskMonitor()

	// Listen for tasks to execute
	go rabbitReceive(SETTINGS.RabbitTaskQueue)

	// Listen for stop commands
	go rabbitReceive(SETTINGS.RabbitStopQueue)

	// cleanly exit after all goroutines are finished
	WG.Wait()
}
