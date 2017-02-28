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
	var t Task = Task{Station: Station{}, StationSchedule: StationSchedule{}}
	var msr ManualStationRun
	ERR = json.NewDecoder(req.Body).Decode(&msr)

	db.Where("id = ?", msr.StationID).Find(&t.Station)
	t.StationSchedule = StationSchedule{Duration: msr.Duration}
	if SETTINGS.PirriDebug {
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
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"stations\": "+string(blob)+"}")
}

func stationGetWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	stationID, err := strconv.Atoi(req.URL.Query()["stationid"][0])

	db.Where("id = ?", stationID).Find(&station)
	blob, err := json.Marshal(&station)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, string(blob))
}
