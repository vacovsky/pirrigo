package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func weatherCurrentWeb(rw http.ResponseWriter, req *http.Request) {
	w := getCurrentWeather()
	blob, err := json.Marshal(w)
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(rw, "{ \"weather\": "+string(blob)+"}")
}
