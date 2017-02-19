package main

import (
	"encoding/json"
	"fmt"
)

func reactToGpioMessage(m []byte) {
	task := Task{}
	json.Unmarshal(m, &task)
	//	WG.Add(1)
	task.Log()
	//	go
	gpioActivator(task.Station.GPIO, true, task.StationSchedule.Duration)
	//	WG.Done()
}

func reactToStopMessage(m []byte) {
	KILL = !KILL
	fmt.Println("Paused:", KILL)
}
