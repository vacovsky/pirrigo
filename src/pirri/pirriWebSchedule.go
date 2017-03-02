package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func stationScheduleAllWeb(rw http.ResponseWriter, req *http.Request) {
	stationSchedules := []StationSchedule{}

	db.Where("end_date > ? AND start_date <= ?", time.Now(), time.Now()).Find(&stationSchedules).Order("ASC")
	blob, err := json.Marshal(&stationSchedules)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"stationSchedules\": "+string(blob)+"}")
}

func stationScheduleEditWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	ERR = json.NewDecoder(req.Body).Decode(&scheduleItem)

	if db.NewRecord(&scheduleItem) {
		db.Create(&scheduleItem)
	} else {
		db.Save(&scheduleItem)
	}
	if SETTINGS.PirriDebug {
		spew.Dump(scheduleItem)
	}
	stationScheduleAllWeb(rw, req)
}

func stationScheduleDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	ERR = json.NewDecoder(req.Body).Decode(&scheduleItem)

	db.Delete(&scheduleItem)
	if SETTINGS.PirriDebug {
		spew.Dump(scheduleItem)
	}
	stationScheduleAllWeb(rw, req)
}
