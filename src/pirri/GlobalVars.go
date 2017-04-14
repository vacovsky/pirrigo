package main

import (
	"sync"
)

//SETTINGS stores global variables (connection string data, ports, etc)
var SETTINGS Settings

//SQLCONNSTRING stores the fully-constructed SQL connection string
var SQLCONNSTRING string

//RMQCONNSTRING stores the fully-constructed RabbitMQ connection string
var RMQCONNSTRING string

//COMMONWIRE is the GPIO pin which is connected to the relay port for the common wire needed to activate the solenoid
var COMMONWIRE int

//VERSION is the version of the application
var VERSION = "0.0.2"

//ERR is a global bucket to hold errors.  Will be going away as the application matures.
var ERR error

//WG is the WaitGroup tracker for the applications GoRoutines
var WG sync.WaitGroup

//RUNSTATUS indicates progress for a running task, and also indicates if the routine is idle
var RUNSTATUS = RunStatus{
	IsIdle: true,
}
