package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	Station         Station         `json:"station`         //`gorm:"ForeignKey:Station"`
	StationSchedule StationSchedule `json:"stationSchedule` //`gorm:"ForeignKey:StationSchedule"`
}

func (t *Task) Log() {
	if SETTINGS.PirriDebug {
		fmt.Println("Logging task", t.Station.ID, t.StationSchedule.StartTime)
	}
	GormDbConnect()
	defer db.Close()
	db.Create(&StationHistory{
		StationID:  t.Station.ID,
		ScheduleID: t.StationSchedule.ID,
		Duration:   t.StationSchedule.Duration,
		StartTime:  time.Now(),
	})
}

func (t *Task) Send() {
	taskBlob, ERR := json.Marshal(&t)
	failOnError(ERR, "Could not jsonify task.")
	if SETTINGS.PirriDebug {
		fmt.Println("Sending Task:", string(taskBlob))
	}
	rabbitSend(SETTINGS.RabbitTaskQueue, string(taskBlob))
}

func (t *Task) Execute() {
	if SETTINGS.PirriDebug {
		fmt.Println("Executing task:", t.Station.ID, t.StationSchedule.StartTime)
	}
	gpioActivator(t.Station.GPIO, true, t.StationSchedule.Duration)
}
