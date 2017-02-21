package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/newrelic/go-agent"
)

func startPirriWebApp() {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/station/run":  stationRunWeb,
		"/station/all":  stationAllWeb,
		"/station":      stationGetWeb,
		"/schedule/all": stationScheduleAllWeb,
		"/history":      historyAllWeb,

		"/static/": (func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		}),

		"/": (func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/index.html")
		}),
	}

	if SETTINGS.UseNewRelic {
		SETTINGS.NewRelicLicense = loadNewRelicKey(SETTINGS.NewRelicLicensePath)
		config := newrelic.NewConfig("PirriGo v"+VERSION, SETTINGS.NewRelicLicense)
		NRAPPMON, ERR := newrelic.NewApplication(config)
		fmt.Println("USing NewRelic Monitoring Agent")
		if NRAPPMON == nil || ERR != nil {
			fmt.Println("Unable to load New Relic Agent using given configuration.")
		} else {
			for k, v := range routes {
				http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, k, v))
			}
		}

	} else {
		for k, v := range routes {
			fmt.Println("Not using New Relic for", k)
			http.HandleFunc(k, v)
		}
	}

	// Host server
	panic(http.ListenAndServe(":"+SETTINGS.HttpPort, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
