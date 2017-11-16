package main

import (
	"bufio"
	"fmt"
	"os"

	"./data"
	"./logging"
	"./settings"
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
	//	firstRunDBSetup()

	// check if we are in local debug mode, or actually doing work.
	// If not debug, reset the GPIO state
	if !set.Debug.SimulateGPIO {
		gpioClear()
	}

	// set the common wire for powering solenoids
	setCommonWire()

	// init waitgroups for concurrent processing
	WG.Add(3)

	// Start the Web application for management of schedule etc.
	go startPirriWebApp()

	// Monitor database for pre-scheduled tasks
	go startTaskMonitor()

	// Listen for tasks to execute
	if SETTINGS.Pirri.UseRabbitMQ {
		go rabbitReceive(SETTINGS.RabbitMQ.TaskQueue)
	} else {
		go listenForTasks()
	}

	go listenForExit()

	WG.Wait()
	fmt.Println("Exit key received - exiting!")
}

func showVersion() {
	name := "PirriGo v" + VERSION
	fmt.Println(name)
}

func listenForExit() {
	fmt.Println("=================== PRESS <ENTER> KEY TO EXIT ===================\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	getLogger().LogEvent("PirriGo v" + VERSION + " exiting due to the exit key being pressed.  You did this...")
	WG.Done()
	os.Exit(0)
}
