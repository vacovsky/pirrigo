package main

import (
	"encoding/json"
	"fmt"
)

func reactToGpioMessage(m []byte) {
	task := Task{}
	json.Unmarshal(m, &task)
	task.log()
	task.execute()
}

func reactToStopMessage(m []byte) {
	KILL = !KILL
	fmt.Println("Paused:", KILL)
}
