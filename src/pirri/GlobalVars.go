package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/newrelic/go-agent"
)

//SETTINGS stores global variables (connection string data, ports, etc)
var SETTINGS Settings

//SQLCONNSTRING stores the fully-constructed SQL connection string
var SQLCONNSTRING string

//RMQCONNSTRING stores the fully-constructed RabbitMQ connection string
var RMQCONNSTRING string

//KILL is going to be a goroutine kill switch, but isn't implemented
var KILL bool

//VERSION is the version of the application
var VERSION = "0.0.1"

//ERR is a global bucket to hold errors.  Will be going away as the application matures.
var ERR error

//WG is the WaitGroup tracker for the applications GoRoutines
var WG sync.WaitGroup

//NRAPPMON is used for NewRelic's newrelic.Application struct
var NRAPPMON newrelic.Application

var RUNSTATUS RunStatus = RunStatus{
	IsIdle: true,
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
