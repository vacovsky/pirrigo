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

func nodeEditWeb(rw http.ResponseWriter, req *http.Request) {
	var node DripNode
	ERR = json.NewDecoder(req.Body).Decode(&node)
	if db.NewRecord(&node) {
		db.Create(&node)
	} else {
		db.Save(&node)
	}
	if SETTINGS.PirriDebug {
		spew.Dump(node)
	}
	nodeAllWeb(rw, req)
}

func nodeUsageStatsWeb(rw http.ResponseWriter, req *http.Request) {
	type waterUsageModel struct {
		RunMins  int
		TotalGPH int
		Notes    string
	}
	results := []waterUsageModel{}

	sqlStr := `
SELECT DISTINCT drip_nodes.station_id,
           SUM((duration / 60 )) as runmins,
           (SELECT sum((gph * count)) as totalgph from drip_nodes where drip_nodes.station_id=station_histories.station_ID) as totalgph,
           stations.notes
       FROM station_histories
       INNER JOIN drip_nodes ON drip_nodes.station_id=station_histories.station_id
       INNER JOIN stations ON stations.id=station_histories.station_id
           WHERE start_time >= (CURRENT_DATE - INTERVAL 30 DAY)
           GROUP BY drip_nodes.station_id
           ORDER BY drip_nodes.station_id ASC;
           `

	db.Raw(sqlStr).Find(&results)
	blob, err := json.Marshal(&results)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"waterUsage\": "+string(blob)+"}")
}

//    sqlStr = """
//        SELECT DISTINCT dripnodes.sid,
//            SUM((duration / 60 )) as runmins,
//            (SELECT sum((gph * count)) as totalgph from dripnodes where dripnodes.sid=history.sid) as totalgph,
//            stations.notes
//        FROM history
//        INNER JOIN dripnodes ON dripnodes.sid=history.sid
//        INNER JOIN stations ON stations.id=history.sid
//            WHERE starttime >= (CURRENT_DATE - INTERVAL 30 DAY)
//            GROUP BY dripnodes.sid
//            ORDER BY dripnodes.sid ASC;
//            """
//    results = {'water_usage': []}
//    for d in sqlConn.read(sqlStr):
//        results['water_usage'].append(
//            {
//                'sid': int(d[0]),
//                'notes': str(d[3]),
//                'run_mins': int(d[1]),
//                'total_gph': float(d[2]),
//                'usage_last_30': float((d[1] / 60) * d[2])
//            }
//        )
