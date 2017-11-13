package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	//	"time"

	"github.com/davecgh/go-spew/spew"
)

func stationRunWeb(rw http.ResponseWriter, req *http.Request) {
	var t = Task{Station: Station{}, StationSchedule: StationSchedule{}}
	var msr ManualStationRun
	err := json.NewDecoder(req.Body).Decode(&msr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("START MSR OBJECT")
	// spew.Dump(msr)
	db.Where("id = ?", msr.StationID).Find(&t.Station)
	t.StationSchedule = StationSchedule{Duration: msr.Duration}
	if SETTINGS.Debug.Pirri {
		spew.Dump(t)
		spew.Dump(msr)
	}
	t.send()
}

func stationAllWeb(rw http.ResponseWriter, req *http.Request) {
	stations := []Station{}

	db.Limit(100).Find(&stations)
	blob, err := json.Marshal(&stations)
	if err != nil {
		getLogger().LogError("Error while marshalling all stations from SQL.", err.Error())
	}
	io.WriteString(rw, "{ \"stations\": "+string(blob)+"}")
}

func stationGetWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	stationID, err := strconv.Atoi(req.URL.Query()["stationid"][0])

	db.Where("id = ?", stationID).Find(&station)
	blob, err := json.Marshal(&station)
	if err != nil {
		getLogger().LogError("Error while marshalling single station from SQL.", err.Error())
	}
	io.WriteString(rw, string(blob))
}

func stationEditWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		getLogger().LogError("Error while editing a station.", err.Error())
	}
	if db.NewRecord(&station) {
		db.Create(&station)
	} else {
		db.Save(&station)
	}
	if SETTINGS.Debug.Pirri {
		spew.Dump(station)
	}
	stationAllWeb(rw, req)
}

func stationAddWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		getLogger().LogError("Error while adding a station.", err.Error())
	}
	db.Create(&station)
	stationAllWeb(rw, req)
}

func stationDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		getLogger().LogError("Error while deleting a station.", err.Error())
	}
	if SETTINGS.Debug.Pirri {
		spew.Dump(&station)
	}
	db.Delete(&station)
	stationAllWeb(rw, req)
}
