package pirri

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
	//	"time"
)

func stationRunWeb(rw http.ResponseWriter, req *http.Request) {
	var t = Task{Station: Station{}, StationSchedule: StationSchedule{}}
	var msr ManualStationRun
	err := json.NewDecoder(req.Body).Decode(&msr)
	if err != nil {
		logging.Service().LogError("Unable to execute station ad hoc task submission.", zap.String("error", err.Error()))
	}
	logging.Service().LogEvent("Run event received from web interface for station",
		zap.Int("stationID", msr.StationID),
		zap.Int("durationSeconds", msr.Duration),
	)
	data.Service().DB.Where("id = ?", msr.StationID).Find(&t.Station)
	t.StationSchedule = StationSchedule{Duration: msr.Duration}
	t.send()
}

func stationAllWeb(rw http.ResponseWriter, req *http.Request) {
	stations := []Station{}

	data.Service().DB.Limit(100).Find(&stations)
	blob, err := json.Marshal(&stations)
	if err != nil {
		logging.Service().LogError("Error while marshalling all stations from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"stations\": "+string(blob)+"}")
}

func stationGetWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	stationID, err := strconv.Atoi(req.URL.Query()["stationid"][0])
	if err != nil {
		logging.Service().LogError("Error while loading a station.",
			zap.String("error", err.Error()))
	}

	data.Service().DB.Where("id = ?", stationID).Find(&station)
	blob, err := json.Marshal(&station)
	if err != nil {
		logging.Service().LogError("Error while marshalling single station from SQL.",
			zap.String("error", err.Error()),
			zap.String("stationID", strconv.Itoa(stationID)),
		)
	}
	io.WriteString(rw, string(blob))
}

func stationEditWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		logging.Service().LogError("Error while editing a station.",
			zap.String("error", err.Error()))
	}
	if data.Service().DB.NewRecord(&station) {
		data.Service().DB.Create(&station)
	} else {
		data.Service().DB.Save(&station)
	}

	stationAllWeb(rw, req)
}

func stationAddWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		logging.Service().LogError("Error while adding a station.", zap.String("error", err.Error()))
	}
	data.Service().DB.Create(&station)
	stationAllWeb(rw, req)
}

func stationDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		logging.Service().LogError("Error while deleting a station.",
			zap.String("error", err.Error()),
			// zap.String()
		)

	}

	data.Service().DB.Delete(&station, station.ID)
	stationAllWeb(rw, req)
}
