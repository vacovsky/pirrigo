package main

import (
	"fmt"
)

var VERSION string = "0.0.1"
var err error

func main() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)

	configInit()
	GormSetup()

	//	fmt.Println(GetCurrentTasks())
}
