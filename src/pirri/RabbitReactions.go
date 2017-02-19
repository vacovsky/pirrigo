package main

import (
	"encoding/json"
	"fmt"
)

func reactToGpioMessage(m []byte) {
	task := Task{}
	json.Unmarshal(m, &task)
	task.Log()
	task.Execute()

}

func reactToStopMessage(m []byte) {
	KILL = !KILL
	fmt.Println("Paused:", KILL)
}
