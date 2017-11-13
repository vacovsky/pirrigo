package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//	"strconv"
	"time"
)

func stationScheduleAllWeb(rw http.ResponseWriter, req *http.Request) {
	stationSchedules := []StationSchedule{}
	db.Where("end_date > ? AND start_date <= ?", time.Now(), time.Now()).Find(&stationSchedules).Order("ASC")
	blob, err := json.Marshal(&stationSchedules)
	if err != nil {
		fmt.Println(err)
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
	stationScheduleAllWeb(rw, req)
}

func stationScheduleDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	err := json.NewDecoder(req.Body).Decode(&scheduleItem)
	if err != nil {
		getLogger().LogError("Problem while attempting to decode request body into a station schedule.", err.Error())
	}
	db.Delete(&scheduleItem)
	stationScheduleAllWeb(rw, req)
}
