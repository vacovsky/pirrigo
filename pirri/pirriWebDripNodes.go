package pirri

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
)

func nodeAllWeb(rw http.ResponseWriter, req *http.Request) {
	nodes := []DripNode{}
	data.Service().DB.Find(&nodes)
	blob, err := json.Marshal(&nodes)
	if err != nil {
		logging.Service().LogError("Error displaying all nodes.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"nodes\": "+string(blob)+"}")
}

func nodeAddWeb(rw http.ResponseWriter, req *http.Request) {
	var node DripNode
	err := json.NewDecoder(req.Body).Decode(&node)
	if err != nil {
		logging.Service().LogError("Could not add a node through the web interface.",
			// zap.String("count", strconv.Itoa(node.Count)),
			// zap.String("gph", fmt.Sprintf("%f", node.GPH)),
			// zap.String("nodeID", strconv.Itoa(node.ID)),
			// zap.String("stationID", strconv.Itoa(node.StationID)),
			zap.String("error", err.Error()))
	}
	data.Service().DB.Create(&node)
	nodeAllWeb(rw, req)
}

func nodeDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var node DripNode
	err := json.NewDecoder(req.Body).Decode(&node)
	if err != nil {
		logging.Service().LogError("Could not delete a node through the web interface.", zap.String("error", err.Error()))
	}
	data.Service().DB.Delete(&node)
	nodeAllWeb(rw, req)
}

func nodeEditWeb(rw http.ResponseWriter, req *http.Request) {
	var node DripNode
	err := json.NewDecoder(req.Body).Decode(&node)
	if err != nil {
		logging.Service().LogError("Could not edit a node through the web interface.", zap.String("error", err.Error()))
	}
	data.Service().DB.Save(&node)
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
	   INNER JOIN stations ON stations.id=station_histories.station_id`

	if os.Getenv("PIRRIGO_DB_TYPE") == "mysql" {
		sqlStr += `
           WHERE start_time >= (CURRENT_DATE - INTERVAL 30 DAY)
           GROUP BY drip_nodes.station_id
           ORDER BY drip_nodes.station_id ASC;
		   `
	} else {
		sqlStr += `
           WHERE start_time >= date('now', '-30 DAYS')
           GROUP BY drip_nodes.station_id
		   ORDER BY drip_nodes.station_id ASC;
		   `
	}
	data.Service().DB.Raw(sqlStr).Find(&results)
	for i, r := range results {
		results[i].Total30Days = float32((r.RunMins / 60) * r.TotalGPH)
	}
	blob, err := json.Marshal(&results)
	if err != nil {
		logging.Service().LogError("Unable to parse node usage stats from SQL.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"waterUsage\": "+string(blob)+"}")
}
