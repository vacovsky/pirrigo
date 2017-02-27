package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//	"github.com/davecgh/go-spew/spew"
)

func historyAllWeb(rw http.ResponseWriter, req *http.Request) {
	history := []StationHistory{}
	defer db.Close()

	gormDbConnect()
	db.Find(&history)
	blob, err := json.Marshal(&history)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"history\": "+string(blob)+"}")
}
