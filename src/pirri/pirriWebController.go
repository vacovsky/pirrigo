package main

import (
	"fmt"
	//	"io"
	"net/http"
	//	"runtime"
)

func startPirriWebApp() {
	//	templatePath := "templates/"

	// Station
	http.HandleFunc("/station/add", stationAdd)
	http.HandleFunc("/station/run", stationRun)
	http.HandleFunc("/station/edit", stationEdit)
	http.HandleFunc("/station", stationGet)

	//	// Schedule
	//	http.HandleFunc("/schedule/add", Home)
	//	http.HandleFunc("/schedule/edit", Home)
	//	http.HandleFunc("/schedule", Home)

	//	// History
	//	http.HandleFunc("/history", Home)
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
	//	http.HandleFunc("/", Home)

	preParseTemplates()
	// Host server
	panic(http.ListenAndServe(":"+SETTINGS.HttpPort, nil))
}

func preParseTemplates() {
	//	tmplPath = TmplBasePath + tmplPath
	//	template.Must(template.ParseFiles(tmplPath))
	//	t, _ = template.ParseFiles("templates/index.html") // Parse template file.
}

func Home(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println(logTraffic())

	//	t.Execute(w, nil)
}

func stationGet(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	//	io.WriteString(w, LoadAvailableColors())
}

func stationAdd(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	//	io.WriteString(w, LoadAvailableColors())
}

func stationEdit(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	//	io.WriteString(w, LoadAvailableColors())
}
