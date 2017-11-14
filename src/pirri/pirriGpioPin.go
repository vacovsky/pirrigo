package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"

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

func gpioActivator(t *Task) {
	t.setStatus(true)
	if SETTINGS.Debug.SimulateGPIO {
		gpioSimulation(t.Station.GPIO, true, t.StationSchedule.Duration)
	} else {
		gpioActivate(t.Station.GPIO, true, t.StationSchedule.Duration)
	}
	t.setStatus(false)
}

func gpioSimulation(gpio int, state bool, seconds int) {
	getLogger().LogEvent(fmt.Sprintf(`GPIO Simulation starting. Time: %s; GPIO: %d, State: %t, Duration: %d`,
		zap.Time("startTime", time.Now()),
		zap.Int("gpio", gpio),
		zap.Bool("state", state),
		zap.Int("durationSeconds", seconds)))
	for seconds > 0 && !RUNSTATUS.Cancel {
		time.Sleep(time.Second)
		seconds--
	}
	getLogger().LogEvent(`GPIO Simulation ending.`,
		zap.Time("endTime", time.Now()),
		zap.Int("gpio", gpio),
		zap.Bool("state", state),
		zap.Int("durationSeconds", seconds))
}

func gpioClear() {
	rpio.Open()
	defer rpio.Close()

	gpios := []GpioPin{}
	sql := "SELECT gpio_pins.* FROM gpio_pins WHERE EXISTS(SELECT 1 FROM stations WHERE stations.gpio=gpio_pins.gpio) OR gpio_pins.common = true"
	db.Raw(sql).Find(&gpios)

	for i := range gpios {
		pin := rpio.Pin(gpios[i].GPIO)
		getLogger().LogEvent("Deactivating GPIO",
			zap.Time("endTime", time.Now()),
			zap.Int("gpio", gpios[i].GPIO),
		)
		pin.High()
	}
}

func gpioActivate(gpio int, state bool, seconds int) {
	defer rpio.Close()
	rpio.Open()
	pin := rpio.Pin(gpio)
	common := rpio.Pin(COMMONWIRE)
	pin.Output()
	common.Output()

	getLogger().LogEvent("Activating GPIOs",
		zap.Int("commonWire", COMMONWIRE),
		zap.Int("gpio", gpio),
	)

	common.Low()
	pin.Low()

	// start countdown
	for seconds > 0 && !RUNSTATUS.Cancel {
		time.Sleep(time.Duration(1) * time.Second)
		seconds--
	}
	getLogger().LogEvent("Deactivating GPIOs",
		zap.Int("commonWire", COMMONWIRE),
		zap.Int("gpio", gpio),
	)
	common.High()
	pin.High()
}
