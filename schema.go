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
	)

	var m pirri.Metadata
	err := d.DB.First(&m).Error
	if err != nil {
		// No metadata record exists yet; this is first run
		firstRunDBSetup()
		return
	}
	if !m.FirstRunComplete {
		firstRunDBSetup()
	}
}

func firstRunDBSetup() {
	log := logging.Service()
	log.LogEvent("Beginning first run DB setup...")
	d := data.Service()

	log.LogEvent("Adding set of valid GPIOs.")
	gpioValues := []int{4, 5, 6, 12, 13, 16, 18, 20, 21, 22, 23, 24, 25, 26, 27}
	for _, gpio := range gpioValues {
		d.DB.Create(&pirri.GpioPin{
			GPIO:   gpio,
			Notes:  "",
			Common: false,
		})
	}

	log.LogEvent("Setting common wire relay pin.")
	d.DB.Model(&pirri.GpioPin{}).Where("gpio = ?", 21).Update("common", true)

	log.LogEvent("Inserting example station.")
	exampleStation := pirri.Station{
		GPIO:    4,
		Notes:   "example station",
		Enabled: true,
	}
	d.DB.Create(&exampleStation)

	location, _ := time.LoadLocation("UTC")

	log.LogEvent("Inserting example schedule.")
	d.DB.Create(&pirri.StationSchedule{
		Duration:  10,
		EndDate:   time.Date(2150, 1, 1, 0, 0, 0, 0, location),
		StartDate: time.Date(2020, 1, 1, 0, 0, 0, 0, location),
		StationID: exampleStation.ID,
		Sunday:    true,
		Monday:    true,
		Tuesday:   true,
		Wednesday: true,
		Thursday:  true,
		Friday:    true,
		Saturday:  true,
		StartTime: 1235,
	})

	log.LogEvent("Marking first run as complete.")
	d.DB.Create(&pirri.Metadata{FirstRunComplete: true})

	log.LogEvent("First run setup complete.")
}
