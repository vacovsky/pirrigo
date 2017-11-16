package main

import (
	"encoding/json"
	"io"
	"net/http"

	"./weather"

	"go.uber.org/zap"
)

func weatherCurrentWeb(rw http.ResponseWriter, req *http.Request) {
	w := weather.Service().Current()
	blob, err := json.Marshal(w)
	if err != nil {
		getLogger().LogError("Unable to get current weather.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"weather\": "+string(blob)+"}")
}
