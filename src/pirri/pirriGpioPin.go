package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

//GpioPin - describes a Raspberry Pi GPIO pin
type GpioPin struct {
	ID     int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO   int `gorm:"not null;unique"`
	Notes  string
	Common bool `sql:"DEFAULT:false" gorm:"not null"`
}

func setCommonWire() {
	var gpio GpioPin
	db.Where("common = true").Limit(1).Find(&gpio)
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

func gpioClear() {
	gpios := []GpioPin{}
	sql := "SELECT gpio_pins.* FROM gpio_pins WHERE EXISTS(SELECT 1 FROM stations WHERE stations.gpio=gpio_pins.gpio) OR gpio_pins.common = true"
	db.Raw(sql).Find(&gpios)

	rpio.Open()
	defer rpio.Close()
	for i := range gpios {
		pin := rpio.Pin(gpios[i].GPIO)
		pin.High()
		pin.Low()
	}
}

func gpioActivate(gpio int, state bool, seconds int) {
	rpio.Open()
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

	gpioClear()
}
