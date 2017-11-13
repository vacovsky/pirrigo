package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/newrelic/go-agent"
)

func startPirriWebApp() {
	routes := map[string]func(http.ResponseWriter, *http.Request){
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

		// authentication
		"/login/verify": loginCheck,

		// root
		"/home": webHome,
	}

	if SETTINGS.NewRelic.Active {
		config := newrelic.NewConfig("PirriGo v"+VERSION, SETTINGS.NewRelic.Key)
		NRAPPMON, err := newrelic.NewApplication(config)
		fmt.Println("Using New Relic Monitoring Agent")
		if NRAPPMON == nil || err != nil {
			getLogger().LogEvent("NewRelic being used.")
		} else {
			for k, v := range routes {
				// wrap each route and function in auth handler and new relic
				http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, k, basicAuth(v)))

			}
			// static content does not require authentication
			http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, "/static/", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, r.URL.Path[1:])
			}))

			// routes to the login page if not authenticated, to the main /home otherwise
			http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, "/", loginAuth))
		}
	} else {
		for k, v := range routes {
			fmt.Println("Not using New Relic for", k)
			// wrap each route and function in auth handler
			http.HandleFunc(k, basicAuth(v))
		}
		// static content does not require authentication
		http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})

		// routes to the login page if not authenticated, to the main /home otherwise
		http.HandleFunc("/login", loginAuth)
	}

	// Host server
	panic(http.ListenAndServe(":"+SETTINGS.Web.Port, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
