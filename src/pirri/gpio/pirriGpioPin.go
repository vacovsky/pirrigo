package pirrigpio

import (
	"time"

	"../data"
	"../logging"
	"../settings"
	"github.com/stianeikeland/go-rpio"
	"go.uber.org/zap"
)

//GpioPin - describes a Raspberry Pi GPIO pin
type GpioPin struct {
	ID     int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO   int `gorm:"not null;unique"`
	Notes  string
	Common bool `sql:"DEFAULT:false" gorm:"not null"`
}

func SetCommonWire() {
	d := data.Service()
	var gpio GpioPin
	d.DB.Where("common = true").Limit(1).Find(&gpio)
	settings.Service().GPIO.CommonWire = gpio.GPIO
}

func gpioActivator(t *Task) {
	set := settings.Service()
	t.setStatus(true)
	if set.Debug.SimulateGPIO {
		gpioSimulation(t.Station.GPIO, true, t.StationSchedule.Duration)
	} else {
		gpioActivate(t.Station.GPIO, true, t.StationSchedule.Duration)
	}
	t.setStatus(false)
}

func gpioSimulation(gpio int, state bool, seconds int) {
	log := logging.Service()

	log.LogEvent(`GPIO Simulation starting.`,
		zap.String("startTimeStamp", time.Now().Format(settings.Service().Pirri.DateFormat)),
		zap.Int("gpio", gpio),
		zap.Bool("state", state),
		zap.Int("durationSeconds", seconds))
	for seconds > 0 && !RUNSTATUS.Cancel {
		time.Sleep(time.Second)
		seconds--
	}
	log.LogEvent(`GPIO Simulation ending.`,
		zap.String("endTimeStamp", time.Now().Format(settings.Service().Pirri.DateFormat)),
		zap.Int("gpio", gpio),
		zap.Bool("state", state),
		zap.Int("durationSeconds", seconds))
}

func GPIOClear() {
	log := logging.Service()

	db := data.Service()

	rpio.Open()
	defer rpio.Close()

	gpios := []GpioPin{}
	sql := "SELECT gpio_pins.* FROM gpio_pins WHERE EXISTS(SELECT 1 FROM stations WHERE stations.gpio=gpio_pins.gpio) OR gpio_pins.common = true"
	db.DB.Raw(sql).Find(&gpios)

	for i := range gpios {
		pin := rpio.Pin(gpios[i].GPIO)
		log.LogEvent("Deactivating GPIO",
			zap.Int("gpio", gpios[i].GPIO),
		)
		pin.High()
	}
}

func gpioActivate(gpio int, state bool, seconds int) {
	log := logging.Service()
	set := settings.Service()
	defer rpio.Close()
	rpio.Open()
	pin := rpio.Pin(gpio)
	common := rpio.Pin(set.GPIO.CommonWire)
	pin.Output()
	common.Output()

	log.LogEvent("Activating GPIOs",
		zap.Int("commonWire", set.GPIO.CommonWire),
		zap.Int("gpio", gpio),
		zap.Int("durationSeconds", seconds),
	)

	common.Low()
	pin.Low()

	// start countdown
	for seconds > 0 && !RUNSTATUS.Cancel {
		time.Sleep(time.Duration(1) * time.Second)
		seconds--
	}
	log.LogEvent("Deactivating GPIOs",
		zap.Int("commonWire", set.GPIO.CommonWire),
		zap.Int("gpio", gpio),
		zap.Int("durationSeconds", seconds),
	)
	common.High()
	pin.High()
}
