package main

import (
	"fmt"

	"github.com/vacovsky/pirrigo/src/data"
	"github.com/vacovsky/pirrigo/src/logging"
	"github.com/vacovsky/pirrigo/src/pirri"
	"github.com/vacovsky/pirrigo/src/settings"
)

func main() {

	// initialize dependencies
	set := settings.Service()
	data.Service()
	log := logging.Service()

	fmt.Printf("\nLaunching PirriGo v%s\n\n", set.Pirri.Version)
	log.LogEvent("PirriGo v" + set.Pirri.Version + " starting up")

	// migrate DB schema and populate with seed data
	// TODO: make this nicer.  Check before running anything.
	firstRunDBSetup()

	// check if we are in local debug mode, or actually doing work.
	// If not debug, reset the GPIO state
	if !set.Debug.SimulateGPIO {
		pirri.GPIOClear()
	}

	// set the common wire for powering solenoids
	pirri.SetCommonWire()

	// init waitgroups for concurrent processing
	pirri.WG.Add(3)

	// Start the Web application for management of schedule etc.
	go pirri.StartPirriWebApp()

	// Monitor database for pre-scheduled tasks
	go pirri.StartTaskMonitor()

	// Listen for tasks to execute
	if set.Pirri.UseRabbitMQ {
		go pirri.RabbitReceive(set.RabbitMQ.TaskQueue)
	} else {
		go pirri.ListenForTasks()
	}

	// go listenForExit()

	pirri.WG.Wait()
	fmt.Println("Exit key received - exiting!")
}
