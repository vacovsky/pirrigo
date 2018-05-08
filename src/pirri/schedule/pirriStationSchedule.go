package pirri

import (
	//	"encoding/json"
	"fmt"
	"time"

	"../data"
	"../logging"
	"../settings"
	"go.uber.org/zap"
)

var lastTriggeredItem string

func checkForTasks() {
	scheds := []StationSchedule{}
	nowTime := time.Now()
	sqlFilter := fmt.Sprintf("(start_date <= NOW() AND end_date > NOW()) AND %s=true AND start_time=%s",
		nowTime.Weekday(),
		fmt.Sprintf("%02d%02d",
			nowTime.Hour(),
			nowTime.Minute()))
	data.Service().DB.Where(sqlFilter).Find(&scheds)
	sendFoundScheduleItems(scheds)
}

func StartTaskMonitor() {
	set := settings.Service()
	logging.Service().LogEvent(`Starting monitoring at interval`,
		zap.Int("interval", set.Pirri.MonitorInterval))
	for {
		checkForTasks()
		time.Sleep(time.Duration(set.Pirri.MonitorInterval) * time.Second)
	}
}

func sendFoundScheduleItems(items []StationSchedule) {
	for i := range items {
		task := Task{StationSchedule: items[i]}
		data.Service().DB.Where(Station{ID: task.StationSchedule.StationID}).Find(&task.Station)
		task.send()
	}
}
