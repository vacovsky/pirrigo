package pirri

import (
	"encoding/json"
	"time"

	"../data"
	"../logging"
	"../settings"
	"go.uber.org/zap"
)

//Task describes a Station activation sent to a RabbitMQ server for processing in serial by the application.
type Task struct {
	Station         Station         `json:"station"`         //`gorm:"ForeignKey:Station"`
	StationSchedule StationSchedule `json:"stationSchedule"` //`gorm:"ForeignKey:StationSchedule"`
}

func (t *Task) log() {
	logging.Service().LogEvent("Logging task for station",
		zap.Int("stationID", t.Station.ID),
		zap.Int("startTime", t.StationSchedule.StartTime),
	)
	if t.Station.GPIO > 0 {
		data.Service().DB.Create(&StationHistory{
			StationID:  t.Station.ID,
			ScheduleID: t.StationSchedule.ID,
			Duration:   t.StationSchedule.Duration,
			StartTime:  time.Now(),
		})
	}
}

func (t *Task) send() {
	if t.Station.GPIO > 0 {
		if settings.Service().Pirri.UseRabbitMQ {
			logging.Service().LogEvent("Queuing Task for GPIO activation in RabbitMQ for station", zap.Int("gpio", t.Station.GPIO))
			taskBlob, err := json.Marshal(&t)
			if err != nil {
				logging.Service().LogError("Could not JSONify task for sending.",
					zap.String("error", err.Error()))
			}
			rabbitSend(settings.Service().RabbitMQ.TaskQueue, string(taskBlob))
		} else {
			ORQMutex.Lock()
			logging.Service().LogEvent("Queuing Task for GPIO activation in OfflineRunQueue for station", zap.Int("gpio", t.Station.GPIO))
			OfflineRunQueue = append(OfflineRunQueue, t)
			ORQMutex.Unlock()
		}
	}
}

func (t *Task) execute() {
	logging.Service().LogEvent("Executing task for station", zap.Int("stationID", t.Station.ID))

	if t.Station.GPIO > 0 {
		t.log()
		gpioActivator(t)
	}
	logging.Service().LogEvent("Task execution complete for station", zap.Int("stationID", t.Station.ID))
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
