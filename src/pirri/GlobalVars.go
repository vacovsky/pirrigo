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
		"/gpio/all":        gpioPinsAllWeb,
		"/gpio/available":  gpioPinsAvailableWeb,
		"/gpio/common":     gpioPinsCommonWeb,
		"/gpio/common/set": gpioPinsCommonSetWeb,

		// charts and reporting
		"/stats/1": statsActivityByStation,
		"/stats/2": statsActivityByDayOfWeek,
		"/stats/3": statsActivityPerStationByDOW,
		"/stats/4": statsStationActivity,

		// run status
		"/status/run":    statusRunWeb,
		"/status/cancel": statusRunCancel,

		// nodes
		"/nodes":        nodeAllWeb,
		"/nodes/add":    nodeAddWeb,
		"/nodes/edit":   nodeEditWeb,
		"/nodes/usage":  nodeUsageStatsWeb,
		"/nodes/delete": nodeDeleteWeb,

		// weather
		"/weather/current": weatherCurrentWeb,

		// station
		"/station/run":    stationRunWeb,
		"/station/all":    stationAllWeb,
		"/station/add":    stationAddWeb,
		"/station/edit":   stationEditWeb,
		"/station/delete": stationDeleteWeb,
		"/station":        stationGetWeb,

		// schedule
		"/schedule/all":    stationScheduleAllWeb,
		"/schedule/edit":   stationScheduleEditWeb,
		"/schedule/delete": stationScheduleDeleteWeb,

		// history
		"/history": historyAllWeb,

		// authentication and login page
		"/login/verify": loginCheck,

		// root
		"/home": webHome,

		// logs
		"/logs/all": logsAllWeb,
	}

	unprotectedRoutes = map[string]func(http.ResponseWriter, *http.Request){
		// login banner and randon metadata available to anyone
		"/meta": metadataWeb,
	}
)
