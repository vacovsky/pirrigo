package main

import (
	"fmt"
	"html/template"
	//	"io"
	"net/http"
	"runtime"
)

func startPirriWebApp() {
	// POSTs
	//	http.HandleFunc("/formtest", test)

	// GETs
	//	http.HandleFunc("/colors", Colors)
	http.HandleFunc("/", Home)

	// STATIC
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// SERVER
	panic(http.ListenAndServe(":"+SETTINGS.HttpPort, nil))
}

func Home(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	t := template.New("templates/index.html")
	t, _ = template.ParseFiles("templates/index.html") // Parse template file.
	t.Execute(w, nil)
}

func Colors(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	//	io.WriteString(w, LoadAvailableColors())
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
