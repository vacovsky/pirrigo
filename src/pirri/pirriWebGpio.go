package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
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

func gpioPinsCommonWeb(rw http.ResponseWriter, req *http.Request) {
	gpio := GpioPin{}
	db.Where("common = ?", true).Limit(1).Find(&gpio)

	blob, err := json.Marshal(&gpio)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"gpio\": "+string(blob)+"}")
}

func gpioPinsCommonSetWeb(rw http.ResponseWriter, req *http.Request) {
	gpio := GpioPin{}
	err := json.NewDecoder(req.Body).Decode(&gpio)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	gpio.Common = true
	spew.Dump(gpio)
	db.Exec("UPDATE `gpio_pins` SET `common` = 'false'")
	db.Exec("UPDATE `gpio_pins` SET `common` = 'true' WHERE (gpio = ?)", gpio.GPIO)
	gpioPinsCommonWeb(rw, req)
}
