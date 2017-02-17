package main

import (
	"fmt"
)

var VERSION string = "0.0.1"

func main() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)

	configInit()

	fmt.Println(GetCurrentTasks())
}
