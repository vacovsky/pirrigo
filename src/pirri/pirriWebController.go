package main

import (
	"net/http"
	"runtime"
)

func startPirriWebApp() {
	//	templatePath := "templates/"

	// Station
	//	http.HandleFunc("/station/add", stationAdd)
	http.HandleFunc("/station/run", stationRunWeb)
	//	http.HandleFunc("/station/edit", stationEdit)
	http.HandleFunc("/station/all", stationAllWeb)
	http.HandleFunc("/station", stationGetWeb)

	//	// Schedule
	//	http.HandleFunc("/schedule/add", Home)
	//	http.HandleFunc("/schedule/edit", Home)
	//	http.HandleFunc("/schedule", Home)

	//	// History
	http.HandleFunc("/history", historyAllWeb)
	//	http.HandleFunc("/history/add", Home)
	//	http.HandleFunc("/history/edit", Home)

	//	// Settings
	//	http.HandleFunc("/settings", Home)
	//	http.HandleFunc("/settings/add", Home)
	//	http.HandleFunc("/settings/edit", Home)

	//	// GPIO
	//	http.HandleFunc("/gpio", Home)
	//	http.HandleFunc("/gpio/add", Home)
	//	http.HandleFunc("/gpio/edit", Home)

	//	// Drip Nodes
	//	http.HandleFunc("/dropnode", Home)
	//	http.HandleFunc("/dropnodesadd", Home)
	//	http.HandleFunc("/dropnode/edit", Home)

	// Static content
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// Home
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	// Host server
	panic(http.ListenAndServe(":"+SETTINGS.HttpPort, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
