package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func statusRunWeb(rw http.ResponseWriter, req *http.Request) {
	blob, err := json.Marshal(&RUNSTATUS)
	if err != nil {
		getLogger().LogError("Error while marshalling Run Status from SQL.", err.Error())
	}
	io.WriteString(rw, string(blob))
}

func statusRunCancel(rw http.ResponseWriter, req *http.Request) {
	blob, err := json.Marshal(&RUNSTATUS)
	if err != nil {
		getLogger().LogError("Error while marshalling run status from SQL.", err.Error())
	}
	RUNSTATUS.Cancel = true
	io.WriteString(rw, string(blob))
}
