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
		"/gpio/all":       gpioPinsAllWeb,
		"/gpio/available": gpioPinsAvailableWeb,

		// charts and reporting
		"/stats/1": statsActivityByStation,
		"/stats/2": statsActivityByDayOfWeek,
		"/stats/3": statsActivityPerStationByDOW,
		"/stats/4": statsStationActivity,

		// run status
		"/status/run":    statusRunWeb,
		"/status/cancel": statusRunCancel,

		// nodes
		"/nodes":       nodeAllWeb,
		"/nodes/add":   nodeAddWeb,
		"/nodes/edit":  nodeEditWeb,
		"/nodes/usage": nodeUsageStatsWeb,

		// weather
		// TODO write a better algorithm for weather handling

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

	if SETTINGS.UseNewRelic {
		SETTINGS.NewRelicLicense = loadNewRelicKey(SETTINGS.NewRelicLicensePath)
		config := newrelic.NewConfig("PirriGo v"+VERSION, SETTINGS.NewRelicLicense)
		NRAPPMON, ERR := newrelic.NewApplication(config)
		fmt.Println("Using New Relic Monitoring Agent")
		if NRAPPMON == nil || ERR != nil {
			fmt.Println("Unable to load New Relic Agent using given configuration.")
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
	panic(http.ListenAndServe(":"+SETTINGS.HTTPPort, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
