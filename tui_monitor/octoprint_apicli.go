package main

import (
	"encoding/json"
	"log"
)

// joe
// terminal_cli
// D420DE95C1EF45F2BFEBA00E8D2FB696
// http://toadstool:2090/api/job

// {"job":{"averagePrintTime":null,"estimatedPrintTime":21254.152278805403,"filament":{"tool0":{"length":15565.6529845576,"volume":37.43978300259524}},"file":{"date":1640395582,"display":"grinch.gcode","name":"grinch.gcode","origin":"local","path":"figurines/grinch.gcode","size":21086951},"lastPrintTime":null,"user":"joe"},"progress":{"completion":89.68933441349581,"filepos":18912746,"printTime":18709,"printTimeLeft":2483,"printTimeLeftOrigin":"genius"},"state":"Printing"}
// GET /api/job HTTP/1.1
// Host: example.com
// X-Api-Key: D420DE95C1EF45F2BFEBA00E8D2FB696

type OctoPrintRunStatus struct {
	Job struct {
		EstimatedPrintTime int64 `json:"estimatedPrintTime"`
		Filament           struct {
			Tool0 struct {
				Length int64   `json:"length"`
				Volume float64 `json:"volume"`
			} `json:"tool0"`
		} `json:"filament"`
		File struct {
			Date   int64  `json:"date"`
			Name   string `json:"name"`
			Origin string `json:"origin"`
			Size   int64  `json:"size"`
		} `json:"file"`
	} `json:"job"`
	Progress struct {
		Completion    float64 `json:"completion"`
		Filepos       int64   `json:"filepos"`
		PrintTime     int64   `json:"printTime"`
		PrintTimeLeft int64   `json:"printTimeLeft"`
	} `json:"progress"`
	State string `json:"state"`
}

func getPrintStatus() OctoPrintRunStatus {
	result := OctoPrintRunStatus{}
	_, body, err := callGetPrinterStatus(`http://toadstool:2090/api/job`)
	if err != nil {
		log.Printf(err.Error())
	}
	json.Unmarshal(body, &result)
	return result
}
