package main

import (
	"encoding/json"
	//	"fmt"
	//	"io"
	//	"io/ioutil"
	//	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func stationRun(rw http.ResponseWriter, req *http.Request) {
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
