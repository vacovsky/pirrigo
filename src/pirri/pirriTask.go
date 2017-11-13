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
	getLogger().LogEvent(fmt.Sprintf("Logging task for station ID: %d at %d", t.Station.ID, t.StationSchedule.StartTime))
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
	if t.Station.GPIO > 0 {
		if SETTINGS.Pirri.UseRabbitMQ {
			getLogger().LogEvent(fmt.Sprintf("Queuing Task for GPIO activation in RabbitMQ for station at GPIO: %d", t.Station.GPIO))
			taskBlob, err := json.Marshal(&t)
			if err != nil {
				getLogger().LogError("Could not JSONify task for sending.", err.Error())
			}
			rabbitSend(SETTINGS.RabbitMQ.TaskQueue, string(taskBlob))
		} else {
			ORQMutex.Lock()
			getLogger().LogEvent(fmt.Sprintf("Queuing Task for GPIO activation in OfflineRunQueue for station at GPIO: %d", t.Station.GPIO))
			OfflineRunQueue = append(OfflineRunQueue, t)
			ORQMutex.Unlock()
		}
	}
}

func (t *Task) execute() {
	getLogger().LogEvent(fmt.Sprintf("Executing task for station ID: %d", t.Station.ID))

	if t.Station.GPIO > 0 {
		t.log()
		gpioActivator(t)
	}
	getLogger().LogEvent(fmt.Sprintf("Task execution complete for station ID: %d", t.Station.ID))
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
