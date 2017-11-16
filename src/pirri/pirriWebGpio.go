package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func gpioPinsAllWeb(rw http.ResponseWriter, req *http.Request) {
	gpios := []GpioPin{}
	db.Limit(100).Find(&gpios)

	blob, err := json.Marshal(&gpios)
	if err != nil {
		log.LogError("Error while marshalling GPIO pins from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"gpios\": "+string(blob)+"}")
}

func gpioPinsAvailableWeb(rw http.ResponseWriter, req *http.Request) {
	gpios := []GpioPin{}
	sql := "SELECT gpio_pins.* FROM gpio_pins WHERE NOT EXISTS(SELECT 1 FROM stations WHERE stations.gpio=gpio_pins.gpio) AND gpio_pins.common = false"
	db.Raw(sql).Find(&gpios)

	blob, err := json.Marshal(&gpios)
	if err != nil {
		log.LogError("Error while marshalling GPIO pins from SQL.",
			zap.String("error", err.Error()),
			zap.String("SQLQuery", sql))
	}
	io.WriteString(rw, "{ \"gpios\": "+string(blob)+"}")
}

func gpioPinsCommonWeb(rw http.ResponseWriter, req *http.Request) {
	gpio := GpioPin{}
	db.Where("common = ?", true).Limit(1).Find(&gpio)

	blob, err := json.Marshal(&gpio)
	if err != nil {
		log.LogError("Error while marshalling GPIO pins from SQL.",
			zap.String("error", err.Error()),
			// zap.String("gpio.GPIO", strconv.Itoa(gpio.GPIO)),
			// zap.String("gpio.ID", strconv.Itoa(gpio.ID)),
			// zap.String("gpio.Notes", gpio.Notes),
		)
	}
	io.WriteString(rw, "{ \"gpio\": "+string(blob)+"}")
}

func gpioPinsCommonSetWeb(rw http.ResponseWriter, req *http.Request) {
	gpio := GpioPin{}
	err := json.NewDecoder(req.Body).Decode(&gpio)
	if err != nil {
		log.LogError("Unable to decode request body when setting common GPIO pin.",
			// zap.String("gpio.GPIO", strconv.Itoa(gpio.GPIO)),
			// zap.String("gpio.ID", strconv.Itoa(gpio.ID)),
			// zap.String("gpio.Notes", gpio.Notes),
			zap.String("error", err.Error()))
	}
	gpio.Common = true
	db.Exec("UPDATE `gpio_pins` SET `common` = false")
	db.Exec("UPDATE `gpio_pins` SET `common` = true WHERE (gpio = ?)", gpio.GPIO)
	gpioPinsCommonWeb(rw, req)
}
