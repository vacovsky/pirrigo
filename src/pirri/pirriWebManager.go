package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func stationRunWeb(rw http.ResponseWriter, req *http.Request) {
	var t Task = Task{Station: Station{}, StationSchedule: StationSchedule{}}
	var msr ManualStationRun
	ERR = json.NewDecoder(req.Body).Decode(&msr)

	GormDbConnect()
	defer db.Close()

	db.Where("id = ?", msr.StationID).Find(&t.Station)
	t.StationSchedule = StationSchedule{Duration: msr.Duration}
	if SETTINGS.PirriDebug {
		spew.Dump(t)
		spew.Dump(msr)
	}
	t.Send()
}

func stationAllWeb(rw http.ResponseWriter, req *http.Request) {
	stations := []Station{}

	GormDbConnect()
	defer db.Close()

	db.Limit(100).Find(&stations)
	blob, err := json.Marshal(&stations)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"stations\": "+string(blob)+"}")
}

func stationGetWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	stationId, err := strconv.Atoi(req.URL.Query()["stationid"][0])

	GormDbConnect()
	defer db.Close()

	db.Where("id = ?", stationId).Find(&station)
	blob, err := json.Marshal(&station)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, string(blob))
}

func historyAllWeb(rw http.ResponseWriter, req *http.Request) {
	history := []StationHistory{}

	GormDbConnect()
	defer db.Close()

	db.Find(&history)
	blob, err := json.Marshal(&history)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"history\": "+string(blob)+"}")
}

func stationScheduleAllWeb(rw http.ResponseWriter, req *http.Request) {
	stationSchedules := []StationSchedule{}

	GormDbConnect()
	defer db.Close()

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

	GormDbConnect()
	defer db.Close()

	if db.NewRecord(&scheduleItem) {
		db.Create(&scheduleItem)
	} else {
		db.Update(&scheduleItem)
	}
	if SETTINGS.PirriDebug {
		spew.Dump(scheduleItem)
	}
}

func stationScheduleDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	ERR = json.NewDecoder(req.Body).Decode(&scheduleItem)

	GormDbConnect()
	defer db.Close()

	db.Delete(&scheduleItem)

	if SETTINGS.PirriDebug {
		spew.Dump(scheduleItem)
	}
}
