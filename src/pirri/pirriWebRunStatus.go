package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func statusRunWeb(rw http.ResponseWriter, req *http.Request) {
	blob, err := json.Marshal(&RUNSTATUS)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, string(blob))
}

func statusRunCancel(rw http.ResponseWriter, req *http.Request) {
	blob, err := json.Marshal(&RUNSTATUS)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	RUNSTATUS.Cancel = true
	io.WriteString(rw, string(blob))
}
