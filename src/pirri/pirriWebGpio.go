package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func gpioPinsAllWeb(rw http.ResponseWriter, req *http.Request) {
	gpios := []GpioPin{}
	db.Limit(100).Find(&gpios)

	blob, err := json.Marshal(&gpios)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"gpios\": "+string(blob)+"}")
}

func gpioPinsAvailableWeb(rw http.ResponseWriter, req *http.Request) {
	gpios := []GpioPin{}
	sql := "SELECT gpio_pins.* FROM gpio_pins WHERE NOT EXISTS(SELECT 1 FROM stations WHERE stations.gpio=gpio_pins.gpio) && gpio_pins.common = false"
	db.Raw(sql).Find(&gpios)

	blob, err := json.Marshal(&gpios)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"gpios\": "+string(blob)+"}")
}
