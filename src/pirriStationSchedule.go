package main

import (
	//	"fmt"
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
