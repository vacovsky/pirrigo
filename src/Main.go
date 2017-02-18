package main

import (
	//	"encoding/json"
	"fmt"
)

var VERSION string = "0.0.1"
var err error

func main() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)

	configInit()
	GormSetup()

	CheckForTask()
	//	GetAllGpio()
	//	fmt.Println(GetCurrentTasks())
	//	CreateNewStationSchedule()
}
