package main

import (
	"fmt"
)

type Task struct {
	Station         Station         //`gorm:"ForeignKey:Station"`
	StationSchedule StationSchedule //`gorm:"ForeignKey:StationSchedule"`
}

func (t *Task) Log() {
	fmt.Println("Logging task", t.Station.ID, t.StationSchedule.StartTime)
	// TODO write to database
}
