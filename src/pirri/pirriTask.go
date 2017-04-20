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
	if SETTINGS.Debug.Pirri {
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
	if ERR != nil {
		fmt.Println(ERR, "Could not JSONify task.")
	}
	if SETTINGS.Debug.Pirri {
		fmt.Println("Queuing Task for GPIO activation cr", t.Station.GPIO)
		spew.Dump(OfflineRunQueue)
	}
	if t.Station.GPIO > 0 {
		if SETTINGS.Pirri.UseRabbitMQ {
			rabbitSend(SETTINGS.RabbitMQ.TaskQueue, string(taskBlob))
		} else {
			ORQMutex.Lock()
			OfflineRunQueue = append(OfflineRunQueue, t)
			spew.Dump(OfflineRunQueue)
			ORQMutex.Unlock()
		}
	}
}

func (t *Task) execute() {
	if SETTINGS.Debug.Pirri {
		fmt.Println("Executing task:", t.Station.ID, t.StationSchedule.StartTime)
		spew.Dump(t)
		spew.Dump(RUNSTATUS)
	}
	if t.Station.GPIO > 0 {
		gpioActivator(t)
	}
	if SETTINGS.Debug.Pirri {
		fmt.Println("Execution of task complete.")
		spew.Dump(RUNSTATUS)
	}
}

func (t *Task) setStatus(active bool) {
	if active {
		manual := t.StationSchedule.ID == 0
		RUNSTATUS = RunStatus{
			Duration:  t.StationSchedule.Duration,
			StationID: t.Station.ID,
			IsIdle:    false,
			StartTime: time.Now(),
			IsManual:  manual,
		}
	} else {
		RUNSTATUS = RunStatus{
			IsIdle: true,
		}
	}
}
