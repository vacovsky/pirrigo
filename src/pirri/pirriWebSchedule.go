package pirri

import (
	"encoding/json"
	"io"
	"net/http"
	//	"strconv"
	"time"

	"github.com/vacovsky/pirrigo/src/data"
	"github.com/vacovsky/pirrigo/src/logging"
	"go.uber.org/zap"
)

func stationScheduleAllWeb(rw http.ResponseWriter, req *http.Request) {
	stationSchedules := []StationSchedule{}
	data.Service().DB.Where("end_date > ? AND start_date <= ?", time.Now(), time.Now()).Find(&stationSchedules).Order("ASC")
	blob, err := json.Marshal(&stationSchedules)
	if err != nil {
		logging.Service().LogError("Unable to retrieve station schedules via web interface.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"stationSchedules\": "+string(blob)+"}")
}

func stationScheduleEditWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	err := json.NewDecoder(req.Body).Decode(&scheduleItem)

	if err != nil {
		logging.Service().LogError("Problem while attempting to decode request body into a station schedule.", zap.String("error", err.Error()))
	}

	if data.Service().DB.NewRecord(&scheduleItem) {
		data.Service().DB.Create(&scheduleItem)
	} else {
		data.Service().DB.Save(&scheduleItem)
	}
	stationScheduleAllWeb(rw, req)
}

func stationScheduleDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	err := json.NewDecoder(req.Body).Decode(&scheduleItem)
	if err != nil {
		logging.Service().LogError("Problem while attempting to decode request body into a station schedule.", zap.String("error", err.Error()))
	}
	data.Service().DB.Delete(&scheduleItem)
	stationScheduleAllWeb(rw, req)
}
