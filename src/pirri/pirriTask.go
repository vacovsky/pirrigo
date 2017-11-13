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
	if SETTINGS.Debug.Pirri {
		fmt.Println("Queuing Task for GPIO activation in RabbitMQ: ", t.Station.GPIO)
		spew.Dump(OfflineRunQueue)
	}
	if t.Station.GPIO > 0 {
		if SETTINGS.Pirri.UseRabbitMQ {
			taskBlob, err := json.Marshal(&t)
			if err != nil {
				getLogger().LogError("Could not JSONify task for sending.", err.Error())
			}
			rabbitSend(SETTINGS.RabbitMQ.TaskQueue, string(taskBlob))
		} else {
			ORQMutex.Lock()
			fmt.Println("Queuing Task for GPIO activation in OfflineRunQueue:", t.Station.GPIO)
			OfflineRunQueue = append(OfflineRunQueue, t)
			spew.Dump(t)
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
		t.log()
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
