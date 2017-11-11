package main

import (
	"sync"
)

var (
	//SETTINGS stores global variables (connection string data, ports, etc) which are poplated at start/reload
	SETTINGS Settings

	//SQLConnString stores the fully-constructed SQL connection string
	SQLConnString string

	//RMQCONNSTRING stores the fully-constructed RabbitMQ connection string
	RMQCONNSTRING string

	//COMMONWIRE is the GPIO pin which is connected to the relay port for the common wire needed to activate the solenoid
	COMMONWIRE int

	//VERSION is the version of the application
	VERSION = "0.3.0"

	//ERR is a global bucket to hold errors.  Will be going away as the application matures.
	ERR error

	//WG is the WaitGroup tracker for the applications GoRoutines
	WG sync.WaitGroup

	//RUNSTATUS indicates progress for a running task, and also indicates if the routine is idle
	RUNSTATUS = RunStatus{
		IsIdle: true,
	}

	// OfflineRunQueue is the task queue for non-rabbit configurations
	OfflineRunQueue = []*Task{}

	// ORQMutex protects the OFFLINE_RUN_QUEUE from race conditions
	ORQMutex = &sync.Mutex{}

	// Log is the struct for logging stuff
	// Logger = logging.MustGetLogger("pirri")
	Logger = logHelper{}
)
