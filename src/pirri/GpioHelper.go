package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func setCommonWire() {
	var gpio GpioPin
	db.Where("common = ?", true).Limit(1).Find(&gpio)
	COMMONWIRE = gpio.GPIO
}

func gpioActivator(gpio int, state bool, seconds int) {
	if SETTINGS.Debug.SimulateGPIO {
		gpioSimulation(gpio, state, seconds)
	} else {
		gpioActivate(gpio, state, seconds)
	}
}

func gpioSimulation(gpio int, state bool, seconds int) {
	fmt.Println("GPIO Simulation starting:",
		"\nTime:", time.Now(),
		"\nGPIO:", gpio,
		"\nDesired State:", state,
		"\nDuration (seconds):", seconds)
	fmt.Println("Active!", time.Now())
	for seconds > 0 && !RUNSTATUS.Cancel {
		time.Sleep(time.Second)
		seconds--
	}
	fmt.Println("Deactivated!", time.Now())
}

func gpioActivate(gpio int, state bool, seconds int) {
	rpio.Open()
	defer rpio.Close()
	pin := rpio.Pin(gpio)
	common := rpio.Pin(COMMONWIRE)
	pin.Output()
	common.Output()

	// activate gpio
	if state {
		common.High()
		pin.High()
	} else {
		common.Low()
		pin.Low()
	}

	// start countdown
	for seconds > 0 && RUNSTATUS.Cancel {
		time.Sleep(time.Duration(1) * time.Second)
		seconds--
	}

	// undo what we just did
	if !state {
		common.Low()
		pin.High()
	} else {
		common.Low()
		pin.Low()
	}
}
