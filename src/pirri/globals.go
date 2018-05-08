package pirri

import (
	"net/http"
	"sync"
)

var (
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

	protectedRoutes = map[string]func(http.ResponseWriter, *http.Request){
		// GPIO Pins
		"/api/gpio/all":        gpioPinsAllWeb,
		"/api/gpio/available":  gpioPinsAvailableWeb,
		"/api/gpio/common":     gpioPinsCommonWeb,
		"/api/gpio/common/set": gpioPinsCommonSetWeb,

		// charts and reporting
		"/api/stats/1": statsActivityByStation,
		"/api/stats/2": statsActivityByDayOfWeek,
		"/api/stats/3": statsActivityPerStationByDOW,
		"/api/stats/4": statsStationActivity,

		// run status
		"/api/status/run":    statusRunWeb,
		"/api/status/cancel": statusRunCancel,

		// nodes
		"/api/nodes":        nodeAllWeb,
		"/api/nodes/add":    nodeAddWeb,
		"/api/nodes/edit":   nodeEditWeb,
		"/api/nodes/usage":  nodeUsageStatsWeb,
		"/api/nodes/delete": nodeDeleteWeb,

		// weather
		"/api/weather/current": weatherCurrentWeb,

		// station
		"/api/station/run":    stationRunWeb,
		"/api/station/all":    stationAllWeb,
		"/api/station/add":    stationAddWeb,
		"/api/station/edit":   stationEditWeb,
		"/api/station/delete": stationDeleteWeb,
		"/api/station":        stationGetWeb,

		// schedule
		"/api/schedule/all":    stationScheduleAllWeb,
		"/api/schedule/edit":   stationScheduleEditWeb,
		"/api/schedule/delete": stationScheduleDeleteWeb,

		// history
		"/api/history": historyAllWeb,

		// authentication and login page
		"/api/login/verify": loginCheck,

		// root
		"/api/home": webHome,

		// logs
		"/api/logs/all": logsAllWeb,
	}

	unprotectedRoutes = map[string]func(http.ResponseWriter, *http.Request){
		// login banner and randon metadata available to anyone
		"/api/meta": metadataWeb,
	}
)
