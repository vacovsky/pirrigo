package main

import (
	//	"bufio"
	//	"encoding/json"
	"fmt"
	"html/template"
	"io"
	//	"log"
	"net/http"
	//	"os"
	//	"reflect"
	"runtime"
	//	"strconv"
	//	"strings"
)

func startPirriWebApp() {
	// INIT LOAD CONFIG
	loadConfig()
	displayConfig()

	// POSTs
	http.HandleFunc("/formtest", test)

	// GETs
	http.HandleFunc("/colors", Colors)
	http.HandleFunc("/roomlighting", RoomLightingLoop)
	http.HandleFunc("/fuse", FusePattern)

	http.HandleFunc("/randompattern", RandomPattern)
	http.HandleFunc("/gradient", GradientPattern)
	http.HandleFunc("/gradientloop", GradientPatternLoop)
	http.HandleFunc("/", Home)

	// STATIC
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// SERVER
	panic(http.ListenAndServe(":"+PORT, nil))
}

func Home(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	t := template.New("templates/index.html")
	t, _ = template.ParseFiles("templates/index.html") // Parse template file.
	t.Execute(w, nil)
}

func Colors(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	io.WriteString(w, LoadAvailableColors())
}

func RandomPattern(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	RedisWriter(Party(LIGHTCOUNT, BRIGHTNESS, 0.04, 1, 1, 1, 1), CHANNEL)
}

func GradientPattern(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	RedisWriter(Gradient(LIGHTCOUNT, BRIGHTNESS, 0.04, 1, 1, 1, 1, 150), CHANNEL)
}

func GradientPatternLoop(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	for i := 0; i < 1000; i++ {
		RedisWriter(Gradient(LIGHTCOUNT, BRIGHTNESS, 0.01, 1, 1, 1, 1, 600), CHANNEL)
	}
}

func FusePattern(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	for i := 0; i < 10; i++ {
		RedisWriter(Fuse(LIGHTCOUNT, BRIGHTNESS, 0.01, 1, 1, 1, 1, 600), CHANNEL)
	}
}

func RoomLightingLoop(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	for i := 0; i < 10; i++ {
		RedisWriter(RoomLighting(), CHANNEL)
	}
}

func roomLightingSlider(w http.ResponseWriter, req *http.Request) {
	fmt.Println(logTraffic())
	for i := 0; i < 250; i++ {
		RedisWriter(RoomLighting(), CHANNEL)
	}
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
