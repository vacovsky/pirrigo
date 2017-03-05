package main

import (
	"time"
)

type RunStatus struct {
	IsIdle     bool
	IsManual   bool
	StartTime  time.Time
	Duration   int
	ScheduleID int
	StationID  int
	Cancel     bool
}
