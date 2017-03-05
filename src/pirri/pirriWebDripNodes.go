package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func nodeAllWeb(rw http.ResponseWriter, req *http.Request) {
	nodes := []DripNode{}
	db.Find(&nodes)
	blob, err := json.Marshal(&nodes)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"nodes\": "+string(blob)+"}")
}

func nodeAddWeb(rw http.ResponseWriter, req *http.Request) {
	var node DripNode
	ERR = json.NewDecoder(req.Body).Decode(&node)
	if SETTINGS.PirriDebug {
		spew.Dump(node)
	}
	db.Create(&node)
	nodeAllWeb(rw, req)
}

func nodeDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var node DripNode
	ERR = json.NewDecoder(req.Body).Decode(&node)
	if SETTINGS.PirriDebug {
		spew.Dump(node)
	}
	db.Delete(&node)
	nodeAllWeb(rw, req)
}

func nodeEditWeb(rw http.ResponseWriter, req *http.Request) {
	var node DripNode
	ERR = json.NewDecoder(req.Body).Decode(&node)
	db.Save(&node)
	if SETTINGS.PirriDebug {
		spew.Dump(node)
	}
	nodeAllWeb(rw, req)
}

func nodeUsageStatsWeb(rw http.ResponseWriter, req *http.Request) {
	type waterUsageModel struct {
		StationID   int
		RunMins     float32
		TotalGPH    float32
		Notes       string
		Total30Days float32
	}
	results := []waterUsageModel{}

	sqlStr := `
SELECT DISTINCT drip_nodes.station_id,
           SUM((duration / 60 )) as run_mins,
           (SELECT sum((gph * count) + 0.0) as total_gph from drip_nodes where drip_nodes.station_id=station_histories.station_ID) as total_gph,
           stations.notes
       FROM station_histories
       INNER JOIN drip_nodes ON drip_nodes.station_id=station_histories.station_id
       INNER JOIN stations ON stations.id=station_histories.station_id
           WHERE start_time >= (CURRENT_DATE - INTERVAL 30 DAY)
           GROUP BY drip_nodes.station_id
           ORDER BY drip_nodes.station_id ASC;
           `
	db.Raw(sqlStr).Find(&results)
	for i, r := range results {
		results[i].Total30Days = float32((r.RunMins / 60) * r.TotalGPH)
	}
	blob, err := json.Marshal(&results)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"waterUsage\": "+string(blob)+"}")
}
