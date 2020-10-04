package pirri

import (
	"time"
)

//RunStatus indicates progress of a task, or if the routine is idle
type RunStatus struct {
	IsIdle     bool
	IsManual   bool
	StartTime  time.Time
	Duration   int
	ScheduleID int
	StationID  int
	Cancel     bool
}
