package pirri

import (
	"encoding/json"
	"io"
	"net/http"
	//	"strconv"
	"time"

	"go.uber.org/zap"
)

func stationScheduleAllWeb(rw http.ResponseWriter, req *http.Request) {
	stationSchedules := []StationSchedule{}
	db.Where("end_date > ? AND start_date <= ?", time.Now(), time.Now()).Find(&stationSchedules).Order("ASC")
	blob, err := json.Marshal(&stationSchedules)
	if err != nil {
		log.LogError("Unable to retrieve station schedules via web interface.", zap.String("error", err.Error()))
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
		log.LogError("Problem while attempting to decode request body into a station schedule.", zap.String("error", err.Error()))
	}
	db.Delete(&scheduleItem)
	stationScheduleAllWeb(rw, req)
}
