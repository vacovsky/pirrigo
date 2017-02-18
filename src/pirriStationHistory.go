package main

import (
	//	"fmt"
	"time"
)

type StationHistory struct {
	ID         int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	SID        int
	ScheduleID int
	Duration   int
	StartTime  time.Time `sql:"DEFAULT:current_timestamp" gorm:"primary_key"`
}
