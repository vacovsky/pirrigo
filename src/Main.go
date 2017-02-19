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

	WG.Add(2)
	go GpioActivator(4, true, 300)
	go TaskMonitor()

	// cleanly exit after all goroutines are finished
	WG.Wait()
}
