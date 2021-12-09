package pirri

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
	//	"time"
)

// TODO parameterize the inputs for date ranges and add selectors on stats page
func statsActivityByStation(rw http.ResponseWriter, req *http.Request) {
	type StatsChart struct {
		ReportType int
		Labels     []int
		Series     []string
		Data       [][]int
	}
	type RawResult struct {
		StationID int
		Secs      int
	}
	query := req.URL.Query()
	st := query.Get("startDate")
	et := query.Get("endDate")
	result := StatsChart{
		ReportType: 1,
		Labels:     []int{},
		Series:     []string{"Unscheduled", "Scheduled"},
	}

	var rawResult0 []RawResult
	var rawResult1 []RawResult
	var sqlQuery0 string
	var sqlQuery1 string

	seriesTracker := map[int]int{}
	tracker0 := 0
	tracker1 := 0

	// unscheduled
	sqlQuery0 = fmt.Sprintf(`SELECT DISTINCT station_id, SUM(duration) as secs
				FROM station_histories
				WHERE start_time >= datetime(%s, 'unixepoch')
				AND start_time <= datetime(%s, 'unixepoch')
				AND schedule_id=0 
				AND station_id > 0
				GROUP BY station_id
				ORDER BY station_id ASC`, st, et)

	// scheduled
	sqlQuery1 = fmt.Sprintf(`SELECT DISTINCT station_id, SUM(duration) as secs
				FROM station_histories
				WHERE start_time >= datetime(%s, 'unixepoch')
				AND start_time <= datetime(%s, 'unixepoch')
				AND schedule_id>=1 
				AND station_id > 0
				GROUP BY station_id
				ORDER BY station_id ASC`, st, et)

	data.Service().DB.Raw(sqlQuery0, 7).Scan(&rawResult0)
	data.Service().DB.Raw(sqlQuery1, 7).Scan(&rawResult1)
	result.Data = [][]int{{}, {}}

	for _, i := range rawResult0 {
		result.Data[0] = append(result.Data[0], 0)
		result.Data[1] = append(result.Data[1], 0)

		if loc, ok := seriesTracker[i.StationID]; ok {
			result.Data[0][loc] += i.Secs / 60
		} else {
			seriesTracker[i.StationID] = tracker0
			result.Labels = append(result.Labels, i.StationID)
			result.Data[0][tracker0] += i.Secs / 60
			tracker0++
		}
	}
	for _, i := range rawResult1 {
		if loc, ok := seriesTracker[i.StationID]; ok {
			result.Data[1][loc] += i.Secs / 60
		} else {
			seriesTracker[i.StationID] = tracker1
			result.Labels = append(result.Labels, i.StationID)
			result.Data[1][tracker1] += i.Secs / 60
			tracker1++
		}
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling usage stats.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func statsActivityByDayOfWeek(rw http.ResponseWriter, req *http.Request) {
	type StatsChart struct {
		ReportType int
		Labels     []string
		Series     []string
		Data       [][]int
	}

	result := StatsChart{
		ReportType: 2,
		Labels:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Series:     []string{"Total", "Scheduled", "Unscheduled"},
	}

	type RawResult struct {
		Day  int
		Secs int
	}

	query := req.URL.Query()
	st := query.Get("startDate")
	et := query.Get("endDate")

	var rawResults0 []RawResult
	var rawResults1 []RawResult
	var rawResults2 []RawResult

	var sqlQuery0 string
	var sqlQuery1 string
	var sqlQuery2 string

	sqlQuery0 = fmt.Sprintf(`SELECT DISTINCT strftime( '%%w', (datetime(start_time, '? HOURS'))) as day, SUM(duration) as secs
            FROM station_histories
			WHERE start_time >= datetime(%s, 'unixepoch')
			AND start_time <= datetime(%s, 'unixepoch')
            GROUP BY day
			ORDER BY day ASC`, st, et)
	sqlQuery1 = fmt.Sprintf(`SELECT DISTINCT strftime( '%%w', (datetime(start_time, '? HOURS'))) as day, SUM(duration) as secs
            FROM station_histories
			WHERE start_time >= datetime(%s, 'unixepoch')
			AND start_time <= datetime(%s, 'unixepoch')
			AND schedule_id > 0
            GROUP BY day
			ORDER BY day ASC`, st, et)
	sqlQuery2 = fmt.Sprintf(`SELECT DISTINCT strftime( '%%w', (datetime(start_time, '? HOURS'))) as day, SUM(duration) as secs
            FROM station_histories
			WHERE start_time >= datetime(%s, 'unixepoch')
			AND start_time <= datetime(%s, 'unixepoch')
			AND schedule_id = 0
            GROUP BY day
			ORDER BY day ASC`, st, et)

	data.Service().DB.Raw(sqlQuery0, os.Getenv("PIRRIGO_UTC_OFFSET"), 7).Scan(&rawResults0)
	data.Service().DB.Raw(sqlQuery1, os.Getenv("PIRRIGO_UTC_OFFSET"), 7).Scan(&rawResults1)
	data.Service().DB.Raw(sqlQuery2, os.Getenv("PIRRIGO_UTC_OFFSET"), 7).Scan(&rawResults2)

	result.Data = [][]int{
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	}

	for _, v := range rawResults0 {
		result.Data[0][v.Day-1] = v.Secs / 60
	}
	for _, v := range rawResults1 {
		result.Data[1][v.Day-1] = v.Secs / 60
	}
	for _, v := range rawResults2 {
		result.Data[2][v.Day-1] = v.Secs / 60
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling usage stats.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func statsActivityPerStationByDOW(rw http.ResponseWriter, req *http.Request) {
	type StatsChart struct {
		ReportType int
		Labels     []string
		Series     []string
		Data       [][]float32
	}

	type RawResult struct {
		Day  int
		Mins float32
	}

	result := StatsChart{
		ReportType: 2,
		Labels:     []string{},
		Series:     []string{},
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling usage stats.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func statsStationActivity(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	st := query.Get("startDate")
	et := query.Get("endDate")

	// DataService().DB.Where("timestamp > ?", start).Where("timestamp <= ?", end).Find(&rawdata).Order("timestamp desc")

	type StatsChart struct {
		ReportType int
		Labels     []string
		Series     []int
		Data       [][]int
	}

	type ChartData struct {
		ID      int
		Hour    int
		RunSecs int
	}

	var chartData []ChartData
	result := StatsChart{
		ReportType: 4,
		Labels: []string{"00:00", "01:00", "02:00", "03:00",
			"04:00", "05:00", "06:00", "07:00", "08:00",
			"09:00", "10:00", "11:00", "12:00", "13:00",
			"14:00", "15:00", "16:00", "17:00", "18:00",
			"19:00", "20:00", "21:00", "22:00", "23:00"},
		Series: []int{},
	}
	result.Data = [][]int{}

	var sqlStr string
	if os.Getenv("PIRRIGO_DB_TYPE") == "mysql" {
		sqlStr = fmt.Sprintf(`SELECT stations.id, 
					  HOUR(start_time + INTERVAL %s HOUR) as hour, 
					  (duration) as run_secs
				FROM station_histories
				JOIN stations ON stations.id = station_histories.station_id
				WHERE start_time >= ? 
					AND stations.id > 0
				ORDER BY station_id ASC`, os.Getenv("PIRRIGO_UTC_OFFSET"))
	} else {
		sqlStr = fmt.Sprintf(`SELECT stations.id, 
			strftime('%%H', time(start_time, '%s HOURS')) as hour, 
			(duration) as run_secs
			FROM station_histories
			JOIN stations ON stations.id = station_histories.station_id
			WHERE start_time >= datetime(%s, 'unixepoch')
				AND start_time <= datetime(%s, 'unixepoch')
				AND stations.id > 0
			ORDER BY station_id ASC`, os.Getenv("PIRRIGO_UTC_OFFSET"), st, et)
	}
	seriesTracker := map[int]int{}

	data.Service().DB.Raw(sqlStr).Scan(&chartData)

	for n, i := range chartData {
		if n == 0 || i.ID != result.Series[len(result.Series)-1] {
			seriesTracker[i.ID] = len(result.Series)
			result.Data = append(result.Data, []int{
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0})
			result.Series = append(result.Series, i.ID)
		}
		result.Data[seriesTracker[i.ID]][i.Hour] += i.RunSecs / 60
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling usage stats.", zap.String("error", err.Error()))
	}

	io.WriteString(rw, string(blob))
}
