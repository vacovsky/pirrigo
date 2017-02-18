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
	startTime, _ := strconv.Atoi(fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))

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
	if ERR != nil {
		panic(ERR.Error())
	}
	JsonifySqlResults(result)
}

func CheckForTask() {
	GormDbConnect()
	defer db.Close()

	sched := StationSchedule{}
	nowTime := time.Now()
	sqlFilter := fmt.Sprintf("(start_date <= NOW() AND end_date > NOW()) AND %s=true AND start_time=%s", nowTime.Weekday(), fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))
	db.Where(sqlFilter).First(&sched)
	blob, ERR := json.Marshal(&sched)
	fmt.Println(string(blob))

	if ERR != nil {
		panic(ERR.Error())
	}
}

func TaskMonitor() {
	fmt.Println("Starting monitoring at %s second interval...", SETTINGS.MonitorInterval)
	for !KILL {
		CheckForTask()
		time.Sleep(time.Duration(SETTINGS.MonitorInterval) * time.Second)
	}
	defer WG.Done()
}
