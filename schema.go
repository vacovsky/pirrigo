package main

import (
	"time"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/pirri"
)

func migrateDataSchema() {
	d := data.Service()
	d.DB.AutoMigrate(
		&pirri.Station{},
		&pirri.Metadata{},
		&pirri.DripNode{},
		&pirri.GpioPin{},
		&pirri.StationHistory{},
		&pirri.StationSchedule{},
		// &settings.Settings{},
	)

	var m pirri.Metadata
	d.DB.Select(&m)
	if !m.FirstRunComplete {
		firstRunDBSetup()
	}
}

func firstRunDBSetup() {
	log := logging.Service()
	log.LogEvent("Beginning first run DB setup...")
	d := data.Service()
	addGPIOs := `INSERT INTO gpio_pins (gpio) VALUES (4),(5),(6),(12),(13),(16),(18),(20),(21),(22),(23),(24),(25),(26),(27);`
	setCommonWire := `UPDATE gpio_pins SET notes='common' WHERE gpio=21;`
	addDays := `INSERT INTO station_schedules
	('sunday', 'monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'station_id', 'start_time', 'duration') 
	VALUES
	('true', 'true', 'true', 'true', 'true', 'true', 'true', 1, 1235, 60);`

	log.LogEvent("Adding set of valid GPIOs.")
	d.DB.Raw(addGPIOs)

	log.LogEvent("Setting common wire relay pin.")
	d.DB.Raw(setCommonWire)

	log.LogEvent("Inserting days of the week to to station_schedules table.")
	d.DB.Raw(addDays)

	d.DB.Save(pirri.Station{
		GPIO:  100,
		ID:    100,
		Notes: "example node",
	})

	location, _ := time.LoadLocation("UTC")

	d.DB.Save(&pirri.StationSchedule{
		Duration:  10,
		EndDate:   time.Date(2150, 1, 1, 1, 1, 1, 1, location),
		StartDate: time.Date(2020, 1, 1, 1, 1, 1, 1, location),
		ID:        100,
		StationID: 100,
	})

	gpios := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28}
	for pin := range gpios {
		d.DB.Create(&pirri.GpioPin{
			GPIO:   pin,
			Notes:  "",
			Common: false,
		})
	}

	log.LogEvent("First run setup complete.")
}
