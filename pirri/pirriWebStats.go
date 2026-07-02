package pirri

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/vacovsky/pirrigo/data"
	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
)

// parseDaysParam reads ?days=N from query string, defaults to defDays.
func parseDaysParam(req *http.Request, defDays int) int {
	s := req.URL.Query().Get("days")
	if s == "" {
		return defDays
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return defDays
	}
	return n
}

func statsActivityByStation(rw http.ResponseWriter, req *http.Request) {
	days := parseDaysParam(req, 7)
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

	result := StatsChart{
		ReportType: 1,
		Labels:     []int{},
		Series:     []string{"Unscheduled", "Scheduled"},
	}

	var rawResult0, rawResult1 []RawResult
	var sqlQuery0, sqlQuery1 string

	seriesTracker := map[int]int{}

	if os.Getenv("PIRRIGO_DB_TYPE") == "mysql" {
		sqlQuery0 = `SELECT DISTINCT station_id, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY) AND schedule_id=0 AND station_id > 0
			GROUP BY station_id
			ORDER BY station_id ASC`
		sqlQuery1 = `SELECT DISTINCT station_id, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY) AND schedule_id>=1 AND station_id > 0
			GROUP BY station_id
			ORDER BY station_id ASC`
	} else {
		sqlQuery0 = `SELECT DISTINCT station_id, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= date('now', '-? DAYS') AND schedule_id=0 AND station_id > 0
			GROUP BY station_id
			ORDER BY station_id ASC`
		sqlQuery1 = `SELECT DISTINCT station_id, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= date('now', '-? DAYS') AND schedule_id>=1 AND station_id > 0
			GROUP BY station_id
			ORDER BY station_id ASC`
	}

	data.Service().DB.Raw(sqlQuery0, days).Scan(&rawResult0)
	data.Service().DB.Raw(sqlQuery1, days).Scan(&rawResult1)
	result.Data = [][]int{[]int{}, []int{}}

	for _, i := range rawResult0 {
		if loc, ok := seriesTracker[i.StationID]; ok {
			result.Data[0][loc] += i.Secs / 60
		} else {
			idx := len(result.Labels)
			seriesTracker[i.StationID] = idx
			result.Labels = append(result.Labels, i.StationID)
			result.Data[0] = append(result.Data[0], i.Secs/60)
			result.Data[1] = append(result.Data[1], 0)
		}
	}
	for _, i := range rawResult1 {
		if loc, ok := seriesTracker[i.StationID]; ok {
			result.Data[1][loc] += i.Secs / 60
		} else {
			idx := len(result.Labels)
			seriesTracker[i.StationID] = idx
			result.Labels = append(result.Labels, i.StationID)
			result.Data[0] = append(result.Data[0], 0)
			result.Data[1] = append(result.Data[1], i.Secs/60)
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
	days := parseDaysParam(req, 7)
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

	var rawResults0, rawResults1, rawResults2 []RawResult

	var sqlQuery0, sqlQuery1, sqlQuery2 string

	if os.Getenv("PIRRIGO_DB_TYPE") == "mysql" {
		sqlQuery0 = `SELECT DISTINCT DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY)
			GROUP BY day
			ORDER BY day ASC`
		sqlQuery1 = `SELECT DISTINCT DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY) AND schedule_id > 0
			GROUP BY day
			ORDER BY day ASC`
		sqlQuery2 = `SELECT DISTINCT DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY) AND schedule_id = 0
			GROUP BY day
			ORDER BY day ASC`
	} else {
		sqlQuery0 = `SELECT DISTINCT strftime('%%w', datetime(start_time, '? HOURS')) as day, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= date('now', '-? DAYS')
			GROUP BY day
			ORDER BY day ASC`
		sqlQuery1 = `SELECT DISTINCT strftime('%%w', datetime(start_time, '? HOURS')) as day, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= date('now', '-? DAYS') AND schedule_id > 0
			GROUP BY day
			ORDER BY day ASC`
		sqlQuery2 = `SELECT DISTINCT strftime('%%w', datetime(start_time, '? HOURS')) as day, SUM(duration) as secs
			FROM station_histories
			WHERE start_time >= date('now', '-? DAYS') AND schedule_id = 0
			GROUP BY day
			ORDER BY day ASC`
	}
	data.Service().DB.Raw(sqlQuery0, os.Getenv("PIRRIGO_UTC_OFFSET"), days).Scan(&rawResults0)
	data.Service().DB.Raw(sqlQuery1, os.Getenv("PIRRIGO_UTC_OFFSET"), days).Scan(&rawResults1)
	data.Service().DB.Raw(sqlQuery2, os.Getenv("PIRRIGO_UTC_OFFSET"), days).Scan(&rawResults2)

	result.Data = [][]int{
		[]int{0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0},
	}

	for _, v := range rawResults0 {
		dayIdx := v.Day - 1
		if os.Getenv("PIRRIGO_DB_TYPE") != "mysql" {
			dayIdx = v.Day
		}
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[0][dayIdx] = v.Secs / 60
		}
	}
	for _, v := range rawResults1 {
		dayIdx := v.Day - 1
		if os.Getenv("PIRRIGO_DB_TYPE") != "mysql" {
			dayIdx = v.Day
		}
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[1][dayIdx] = v.Secs / 60
		}
	}
	for _, v := range rawResults2 {
		dayIdx := v.Day - 1
		if os.Getenv("PIRRIGO_DB_TYPE") != "mysql" {
			dayIdx = v.Day
		}
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[2][dayIdx] = v.Secs / 60
		}
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling usage stats.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

// statsActivityPerStationByDOW: per-station watering minutes by day of week.
func statsActivityPerStationByDOW(rw http.ResponseWriter, req *http.Request) {
	days := parseDaysParam(req, 7)
	type StatsChart struct {
		ReportType int
		Labels     []string
		Series     []int
		Data       [][]int
	}

	type RawResult struct {
		StationID int
		Day       int
		Mins      int
	}

	result := StatsChart{
		ReportType: 3,
		Labels:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Series:     []int{},
	}

	var rawResults []RawResult
	sqlQuery := `SELECT station_id, strftime('%%w', datetime(start_time, '? HOURS')) as day, SUM(duration)/60 as mins
		FROM station_histories
		WHERE start_time >= date('now', '-? DAYS') AND station_id > 0
		GROUP BY station_id, day
		ORDER BY station_id, day ASC`

	if os.Getenv("PIRRIGO_DB_TYPE") == "mysql" {
		sqlQuery = `SELECT station_id, DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration)/60 as mins
			FROM station_histories
			WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY) AND station_id > 0
			GROUP BY station_id, day
			ORDER BY station_id, day ASC`
	}

	data.Service().DB.Raw(sqlQuery, os.Getenv("PIRRIGO_UTC_OFFSET"), days).Scan(&rawResults)

	seriesTracker := map[int]int{}
	for _, v := range rawResults {
		if _, ok := seriesTracker[v.StationID]; !ok {
			seriesTracker[v.StationID] = len(result.Series)
			result.Series = append(result.Series, v.StationID)
			result.Data = append(result.Data, []int{0, 0, 0, 0, 0, 0, 0})
		}
		dayIdx := v.Day - 1
		if os.Getenv("PIRRIGO_DB_TYPE") != "mysql" {
			dayIdx = v.Day
		}
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[seriesTracker[v.StationID]][dayIdx] = v.Mins
		}
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling usage stats.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func statsStationActivity(rw http.ResponseWriter, req *http.Request) {
	days := parseDaysParam(req, 7)
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
					WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY)
						AND stations.id > 0
					ORDER BY station_id ASC`, os.Getenv("PIRRIGO_UTC_OFFSET"))
		data.Service().DB.Raw(sqlStr, days).Scan(&chartData)
	} else {
		sqlStr = fmt.Sprintf(`SELECT stations.id,
			strftime('%%H', time(start_time, '%s HOURS')) as hour,
			(duration) as run_secs
		FROM station_histories
		JOIN stations ON stations.id = station_histories.station_id
		WHERE start_time >= date('now', '-? DAYS')
			AND stations.id > 0
		ORDER BY station_id ASC`, os.Getenv("PIRRIGO_UTC_OFFSET"))
		data.Service().DB.Raw(sqlStr, days).Scan(&chartData)
	}
	seriesTracker := map[int]int{}

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
