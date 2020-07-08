package pirri

import (
	//	"encoding/json"
	"fmt"
	"time"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/settings"
	"go.uber.org/zap"
)

var lastTriggeredItem string

//StationSchedule describes a scheduled activation for a Station
type StationSchedule struct {
	ID        int       `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	StartDate time.Time `sql:"DEFAULT:current_timestamp" gorm:"not null"`
	EndDate   time.Time `sql:"DEFAULT:'2025-01-01 00:00:00'" gorm:"not null"`
	Sunday    bool      `sql:"DEFAULT:false" gorm:"not null"`
	Monday    bool      `sql:"DEFAULT:false" gorm:"not null"`
	Tuesday   bool      `sql:"DEFAULT:false" gorm:"not null"`
	Wednesday bool      `sql:"DEFAULT:false" gorm:"not null"`
	Thursday  bool      `sql:"DEFAULT:false" gorm:"not null"`
	Friday    bool      `sql:"DEFAULT:false" gorm:"not null"`
	Saturday  bool      `sql:"DEFAULT:false" gorm:"not null"`
	StationID int       `gorm:"not null"`
	StartTime int       `gorm:"not null"`
	Duration  int       `sql:"DEFAULT:0" gorm:"not null"`
	Repeating bool      `sql:"DEFAULT:false" gorm:"not null"`
}

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
