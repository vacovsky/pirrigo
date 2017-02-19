package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

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

func CreateNewStationSchedule() {
	nowTime := time.Now()
	startTime, ERR := strconv.Atoi(fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))
	failOnError(ERR, "Unable to create new station schedule entry.")
	sched := &StationSchedule{
		StationID: 23,
		StartTime: startTime,
		Duration:  300,
	}
	GormDbConnect()
	defer db.Close()

	db.Create(&sched)
}

func CheckForTaskRaw() {
	sched := StationSchedule{}
	nowTime := time.Now()

	date, now := nowTime.Weekday(), fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute())
	fmt.Println(date, now)
	GormDbConnect()
	defer db.Close()
	sqlQuery := fmt.Sprintf(`SELECT * FROM station_schedules
                        WHERE (start_date <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
                                AND end_date > CAST(replace(date(NOW()), '-', '') AS UNSIGNED))
                            #and %s=1
                            and start_time=%s`,
		nowTime.Weekday(),
		fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))
	result := db.Raw(sqlQuery)
	result.Scan(&sched)

	blob, ERR := json.Marshal(&sched)
	fmt.Println(string(blob))
	failOnError(ERR, "Could not jsonify task.")
	JsonifySqlResults(result)
}

func CheckForTasks() {
	GormDbConnect()
	defer db.Close()

	scheds := []StationSchedule{}
	nowTime := time.Now()
	sqlFilter := fmt.Sprintf("(start_date <= NOW() AND end_date > NOW()) AND %s=true AND start_time=%s", nowTime.Weekday(), fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))

	db.Where(sqlFilter).Find(&scheds)

	if ERR != nil {
		panic(ERR.Error())
	}
}

func TaskMonitor() {
	fmt.Println("Starting monitoring at interval:", SETTINGS.MonitorInterval, "seconds.")
	for !KILL {
		CheckForTasks()
		time.Sleep(time.Duration(SETTINGS.MonitorInterval) * time.Second)
	}
	defer WG.Done()
}

func SendFoundScheduleItems(items []StationSchedule) {
	blob, ERR := json.Marshal(&items)
	failOnError(ERR, "Could not jsonify task.")

	fmt.Println(string(blob))

	for i := range items {
		//		task := Task{}
		fmt.Println(string(i))
	}

}
