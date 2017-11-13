package main

import (
	"encoding/json"
	"io"
	"net/http"
	//	"github.com/davecgh/go-spew/spew"
)

func historyAllWeb(rw http.ResponseWriter, req *http.Request) {
	history := []StationHistory{}

	db.Order("id desc").Limit(100).Find(&history)
	blob, err := json.Marshal(&history)
	if err != nil {
		getLogger().LogError("Error while marshalling history from SQL.", err.Error())
	}
	io.WriteString(rw, "{ \"history\": "+string(blob)+"}")
}
