package main

import (
	"encoding/json"
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

func (t *Task) Send() {
	taskBlob, ERR := json.Marshal(&t)
	failOnError(ERR, "Could not jsonify task.")

	if SETTINGS.PirriDebug {
		fmt.Println("Sending Task:", string(taskBlob))
	}

	rabbitSend(SETTINGS.RabbitTaskQueue, string(taskBlob))
}

func (t *Task) Execute() {
	fmt.Println("Executing task", t.Station.ID, t.StationSchedule.StartTime)
	// TODO Actually execute the task
}
