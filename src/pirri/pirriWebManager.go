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
		"/stats/1": statsActivityByHour,
		"/stats/2": statsActivityByDayOfWeek,
		"/stats/3": statsActivityPerStationByDOW,
		"/stats/4": statsStationActivity,

		// nodes
		// TODO bleh

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

		// static
		"/static/": (func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		}),

		// root
		"/": (func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/index.html")
		}),
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
				http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, k, basicAuth(v)))
			}
		}

	} else {
		for k, v := range routes {
			fmt.Println("Not using New Relic for", k)
			http.HandleFunc(k, basicAuth(v))
		}
	}

	// http.HandleFunc("/", BasicAuth(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "templates/index.html")
	// }, "admin", "123456", "Please enter your username and password for this site"))

	// Host server
	panic(http.ListenAndServe(":"+SETTINGS.HTTPPort, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
