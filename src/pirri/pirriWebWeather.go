package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func weatherCurrentWeb(rw http.ResponseWriter, req *http.Request) {
	w := getCurrentWeather()
	blob, err := json.Marshal(w)
	if err != nil {
		getLogger().LogError("Unable to get current weather.", err.Error())
	}
	io.WriteString(rw, "{ \"weather\": "+string(blob)+"}")
}
