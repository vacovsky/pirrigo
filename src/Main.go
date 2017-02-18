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

	go GpioActivator(4, true, 300)
	go TaskMonitor()
}
