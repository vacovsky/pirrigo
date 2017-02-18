package main

import (
	"fmt"
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
	Fridayb   bool      `sql:"DEFAULT:false"`
	Saturday  bool      `sql:"DEFAULT:false"`
	SID       int
	StartTime int
	Duration  int  `sql:"DEFAULT:0"`
	Repeating bool `sql:"DEFAULT:false"`
}

func CheckForTask() {
	sched := StationSchedule{}
	nowTime := time.Now()

	date, now := nowTime.Weekday(), fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute())
	fmt.Println(date, now)
	GormDbConnect()
	defer db.Close()
	sqlQuery := fmt.Sprintf(`SELECT * FROM schedule
                        WHERE (startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
                                AND enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED))
                            and %s=1
                            and starttime=%s`,
		nowTime.Weekday(),
		fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))
	result := db.Raw(sqlQuery).Scan(&sched)
	fmt.Println(sqlQuery)
	JsonifyResults(result)
}

//	db.Select(`SELECT id, station, duration FROM schedule
//                        WHERE (startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
//                                AND enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED)))`)
