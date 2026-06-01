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
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}
	logging.Service().LogEvent("Run event received from web interface for station",
		zap.Int("stationID", msr.StationID),
		zap.Int("durationSeconds", msr.Duration),
	)
	data.Service().DB.Where("id = ?", msr.StationID).Find(&t.Station)
	t.StationSchedule = StationSchedule{Duration: msr.Duration}
	t.send()
	stationAllWeb(rw, req)
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
	stationIDStr := req.URL.Query().Get("stationid")
	if stationIDStr == "" {
		http.Error(rw, "Missing stationid parameter", http.StatusBadRequest)
		return
	}
	stationID, err := strconv.Atoi(stationIDStr)
	if err != nil {
		logging.Service().LogError("Error while loading a station.",
			zap.String("error", err.Error()))
		http.Error(rw, "Invalid stationid parameter", http.StatusBadRequest)
		return
	}

	data.Service().DB.Where("id = ?", stationID).Find(&station)
	blob, err := json.Marshal(&station)
	if err != nil {
		logging.Service().LogError("Error while marshalling single station from SQL.",
			zap.String("error", err.Error()),
			zap.String("stationID", strconv.Itoa(stationID)),
		)
		http.Error(rw, "Error marshalling station", http.StatusInternalServerError)
		return
	}
	io.WriteString(rw, string(blob))
}

func stationEditWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		logging.Service().LogError("Error while editing a station.",
			zap.String("error", err.Error()))
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
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
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}
	data.Service().DB.Create(&station)
	stationAllWeb(rw, req)
}

func stationDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	err := json.NewDecoder(req.Body).Decode(&station)
	if err != nil {
		logging.Service().LogError("Error while deleting a station.",
			zap.String("error", err.Error()))
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	data.Service().DB.Delete(&station, station.ID)
	stationAllWeb(rw, req)
}
