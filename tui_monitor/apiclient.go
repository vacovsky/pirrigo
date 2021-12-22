package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/vacovsky/pirrigo/pirri"
)

var (
	PIRRI_HOST = "http://192.168.111.130"
	GH_HOST    = "http://192.168.111.23"
	BT_HOST    = "http://192.168.111.164"
)

func makeGetCall(url string) (*http.Response, []byte, error) {
	hc := http.Client{Timeout: 10 * time.Second}
	r, err := hc.Get(url)
	if err != nil {
		log.Panicln(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}
	return r, body, err
}

func getDataForChart(host string) []float64 {
	startTime, endTime := time.Now().Unix()-21600, time.Now().Unix()
	_, data, err := makeGetCall(fmt.Sprintf("%s/tempchartdata?startTime=%d&endTime=%d", host, startTime, endTime))
	var result Chart
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(data, &result)
	return convertStrToFloats(result.Data[0])
}

func calcTimeDiff(status pirri.RunStatus) (int64, int) {
	remaining := (status.StartTime.Unix() + int64(status.Duration)) - time.Now().Unix()
	perc := (int(remaining) / status.Duration)
	return remaining, perc
}

func getPirriRunStatus() pirri.RunStatus {
	_, data, err := makeGetCall(fmt.Sprintf("%s/status/run", PIRRI_HOST))
	var result pirri.RunStatus
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(data, &result)
	return result
}

func cancelStationRun() {
	//  http://192.168.111.130/status/cancel
}
