package main

import (
	"io"
	"net/http"
)

func logsAllWeb(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "{ \"logs\": []}")
}
