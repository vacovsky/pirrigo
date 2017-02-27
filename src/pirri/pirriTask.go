package main

import (
	"encoding/json"
	"fmt"
	"time"
)

//Task describes a Station activation sent to a RabbitMQ server for processing in serial by the application.
type Task struct {
	Station         Station         `json:"station`         //`gorm:"ForeignKey:Station"`
	StationSchedule StationSchedule `json:"stationSchedule` //`gorm:"ForeignKey:StationSchedule"`
}

func (t *Task) log() {
	if SETTINGS.PirriDebug {
		fmt.Println("Logging task", t.Station.ID, t.StationSchedule.StartTime)
	}
	defer db.Close()
	gormDbConnect()
	db.Create(&StationHistory{
		StationID:  t.Station.ID,
		ScheduleID: t.StationSchedule.ID,
		Duration:   t.StationSchedule.Duration,
		StartTime:  time.Now(),
	})
}

func (t *Task) send() {
	taskBlob, ERR := json.Marshal(&t)
	failOnError(ERR, "Could not jsonify task.")
	if SETTINGS.PirriDebug {
		fmt.Println("Sending Task:", string(taskBlob))
	}
	rabbitSend(SETTINGS.RabbitTaskQueue, string(taskBlob))
}

func (t *Task) execute() {
	if SETTINGS.PirriDebug {
		fmt.Println("Executing task:", t.Station.ID, t.StationSchedule.StartTime)
	}
	gpioActivator(t.Station.GPIO, true, t.StationSchedule.Duration)
}
