//go:build !linux || !arm
// +build !linux !arm

package pirri

import (
	"os"
	"strings"
	"time"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/settings"
	"go.uber.org/zap"
)

// GpioPin - describes a Raspberry Pi GPIO pin
type GpioPin struct {
	ID     int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO   int `gorm:"not null;unique"`
	Notes  string
	Common bool `sql:"DEFAULT:false" gorm:"not null"`
}

func SetCommonWire() {
	d := data.Service()
	var gpio GpioPin
	result := d.DB.Where("common = true").Limit(1).Find(&gpio)
	if result.Error != nil {
		logging.Service().LogError("Error querying common wire GPIO", zap.Error(result.Error))
		return
	}
	if gpio.GPIO == 0 {
		logging.Service().LogError("No common wire GPIO pin configured")
		return
	}
	settings.Service().GPIO.CommonWire = gpio.GPIO
}

func gpioActivator(t *Task) {
	t.setStatus(true)
	if strings.ToLower(os.Getenv("PIRRIGO_SIMULATE_GPIO")) == "true" {
		gpioSimulation(t.Station.GPIO, true, t.StationSchedule.Duration)
	} else {
		gpioActivate(t.Station.GPIO, true, t.StationSchedule.Duration)
	}
	t.setStatus(false)
}

func gpioSimulation(gpio int, state bool, seconds int) {
	log := logging.Service()

	log.LogEvent(`GPIO Simulation starting.`,
		zap.String("startTimeStamp", time.Now().Format(os.Getenv("PIRRIGO_DATE_FORMAT"))),
		zap.Int("gpio", gpio),
		zap.Bool("state", state),
		zap.Int("durationSeconds", seconds))
	for seconds > 0 && !RUNSTATUS.Cancel {
		time.Sleep(time.Second)
		seconds--
	}
	log.LogEvent(`GPIO Simulation ending.`,
		zap.String("endTimeStamp", time.Now().Format(os.Getenv("PIRRIGO_DATE_FORMAT"))),
		zap.Int("gpio", gpio),
		zap.Bool("state", state),
		zap.Int("durationSeconds", seconds))
}

func GPIOClear() {
	log := logging.Service()
	log.LogEvent("GPIO Clear skipped (not on Raspberry Pi)")
}

func gpioActivate(gpio int, state bool, seconds int) {
	log := logging.Service()
	set := settings.Service()

	log.LogEvent("GPIO Activate simulated",
		zap.Int("commonWire", set.GPIO.CommonWire),
		zap.Int("gpio", gpio),
		zap.Int("durationSeconds", seconds),
	)

	// start countdown
	for seconds > 0 && !RUNSTATUS.Cancel {
		time.Sleep(time.Duration(1) * time.Second)
		seconds--
	}
	log.LogEvent("GPIO Deactivate simulated",
		zap.Int("commonWire", set.GPIO.CommonWire),
		zap.Int("gpio", gpio),
		zap.Int("durationSeconds", seconds),
	)
}
