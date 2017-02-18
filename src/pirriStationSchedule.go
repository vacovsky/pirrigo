package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type StationSchedule struct {
	ID        int       `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	StartDate time.Time `sql:"DEFAULT:current_timestamp"`
	EndDate   time.Time `sql:"DEFAULT:'2025-01-01 00:00:00'"`
	Sunday    bool      `sql:"DEFAULT:false"`
	Monday    bool      `sql:"DEFAULT:false"`
	Tuesday   bool      `sql:"DEFAULT:false"`
	Wednesday bool      `sql:"DEFAULT:false"`
	Thursday  bool      `sql:"DEFAULT:false"`
	Friday    bool      `sql:"DEFAULT:false"`
	Saturday  bool      `sql:"DEFAULT:false"`
	StationID int
	StartTime int
	Duration  int  `sql:"DEFAULT:0"`
	Repeating bool `sql:"DEFAULT:false"`
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

	blob, err := json.Marshal(&sched)
	fmt.Println(string(blob))
	if err != nil {
		panic("No results or broken model.")
	}
	JsonifySqlResults(result)
}

func CheckForTask() {
	GormDbConnect()
	defer db.Close()

	sched := StationSchedule{}
	//	nowTime := time.Now()
	//	fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute())
	db.First(&sched)
	blob, err := json.Marshal(&sched)
	fmt.Println(string(blob))
	if err != nil {
		panic("No results or broken model.")
	}
}
