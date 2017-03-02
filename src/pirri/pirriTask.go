package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	if t.Station.GPIO > 0 {
		db.Create(&StationHistory{
			StationID:  t.Station.ID,
			ScheduleID: t.StationSchedule.ID,
			Duration:   t.StationSchedule.Duration,
			StartTime:  time.Now(),
		})
	}
}

func (t *Task) send() {
	taskBlob, ERR := json.Marshal(&t)
	failOnError(ERR, "Could not jsonify task.")
	if SETTINGS.PirriDebug {
		fmt.Println("Sending Task:", string(taskBlob))
		spew.Dump(t)
	}
	if t.Station.GPIO > 0 {
		rabbitSend(SETTINGS.RabbitTaskQueue, string(taskBlob))
	}
}

func (t *Task) execute() {
	if SETTINGS.PirriDebug {
		fmt.Println("Executing task:", t.Station.ID, t.StationSchedule.StartTime)
		spew.Dump(t)
	}
	if t.Station.GPIO > 0 {
		gpioActivator(t.Station.GPIO, true, t.StationSchedule.Duration)
	}
}
