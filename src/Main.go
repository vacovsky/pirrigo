package main

import (
	"encoding/json"
	"fmt"
)

func showVersion() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)
}

func main() {
	configInit()
	GormSetup()

	t := Task{StationID: 23, ScheduleID: 4, Duration: 15}
	task, _ := json.Marshal(t)
	//	failOnError(ERR, ERR.Error())

	RabbitSend(SETTINGS.RabbitTaskQueue, string(task))

	WG.Add(4)
	go GpioActivator(4, true, 300)
	go TaskMonitor()
	go RabbitReceive(SETTINGS.RabbitTaskQueue)
	go RabbitReceive(SETTINGS.RabbitStopQueue)

	// cleanly exit after all goroutines are finished
	WG.Wait()
}
