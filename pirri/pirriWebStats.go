package pirri

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

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

// ponytail: helper to build SQLite date filter — avoids ? inside SQL string literals
// which modernc.org/sqlite can't bind as function arguments
func sqliteDateFilter(days int) string {
	return fmt.Sprintf("date('now', '-%d DAYS')", days)
}

// ponytail: extract date from Go-style timestamp (has fractional seconds + tz offset)
// modernc.org/sqlite strftime() can't parse these, so we do it in Go
func extractDate(ts string) string {
	// Take first 10 chars: "2026-07-02 12:35:44.06580017-07:00" -> "2026-07-02"
	if len(ts) >= 10 {
		return ts[:10]
	}
	return ts
}

// ponytail: extract hour from Go-style timestamp string
func extractHour(ts string) int {
	// "2026-07-02 12:35:44..." -> hour at index 11-12
	if len(ts) >= 13 {
		h, err := strconv.Atoi(ts[11:13])
		if err == nil {
			return h
		}
	}
	return 0
}

// ponytail: extract day-of-week from Go-style timestamp (0=Sunday)
func extractDayOfWeek(ts string) int {
	// Parse as Go time format
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		// Fallback: try without timezone
		t, err = time.Parse("2006-01-02 15:04:05", ts[:19])
		if err != nil {
			return 0
		}
	}
	return int(t.Weekday())
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
	cutOff := sqliteDateFilter(days)

	sqlQuery0 := fmt.Sprintf(`SELECT DISTINCT station_id, SUM(duration) as secs
		FROM station_histories
		WHERE start_time >= %s AND schedule_id=0 AND station_id > 0
		GROUP BY station_id
		ORDER BY station_id ASC`, cutOff)
	sqlQuery1 := fmt.Sprintf(`SELECT DISTINCT station_id, SUM(duration) as secs
		FROM station_histories
		WHERE start_time >= %s AND schedule_id>=1 AND station_id > 0
		GROUP BY station_id
		ORDER BY station_id ASC`, cutOff)

	data.Service().DB.Raw(sqlQuery0).Scan(&rawResult0)
	data.Service().DB.Raw(sqlQuery1).Scan(&rawResult1)
	result.Data = [][]int{[]int{}, []int{}}

	seriesTracker := map[int]int{}
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

	// ponytail: fetch raw timestamps, parse day-of-week in Go (SQlite strftime can't parse Go timestamps)
	type RawResult struct {
		StartTime string
		Secs      int
	}

	var rawResults0, rawResults1, rawResults2 []RawResult
	cutOff := sqliteDateFilter(days)

	sqlQuery0 := fmt.Sprintf(`SELECT start_time, SUM(duration) as secs
		FROM station_histories
		WHERE start_time >= %s
		GROUP BY start_time
		ORDER BY start_time ASC`, cutOff)
	sqlQuery1 := fmt.Sprintf(`SELECT start_time, SUM(duration) as secs
		FROM station_histories
		WHERE start_time >= %s AND schedule_id > 0
		GROUP BY start_time
		ORDER BY start_time ASC`, cutOff)
	sqlQuery2 := fmt.Sprintf(`SELECT start_time, SUM(duration) as secs
		FROM station_histories
		WHERE start_time >= %s AND schedule_id = 0
		GROUP BY start_time
		ORDER BY start_time ASC`, cutOff)

	data.Service().DB.Raw(sqlQuery0).Scan(&rawResults0)
	data.Service().DB.Raw(sqlQuery1).Scan(&rawResults1)
	data.Service().DB.Raw(sqlQuery2).Scan(&rawResults2)

	result.Data = [][]int{
		[]int{0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0},
	}

	for _, v := range rawResults0 {
		dayIdx := extractDayOfWeek(v.StartTime)
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[0][dayIdx] += v.Secs / 60
		}
	}
	for _, v := range rawResults1 {
		dayIdx := extractDayOfWeek(v.StartTime)
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[1][dayIdx] += v.Secs / 60
		}
	}
	for _, v := range rawResults2 {
		dayIdx := extractDayOfWeek(v.StartTime)
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[2][dayIdx] += v.Secs / 60
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
		StartTime string
		Secs      int
	}

	result := StatsChart{
		ReportType: 3,
		Labels:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Series:     []int{},
	}

	var rawResults []RawResult
	cutOff := sqliteDateFilter(days)

	sqlQuery := fmt.Sprintf(`SELECT station_id, start_time, SUM(duration) as secs
		FROM station_histories
		WHERE start_time >= %s AND station_id > 0
		GROUP BY station_id, start_time
		ORDER BY station_id, start_time ASC`, cutOff)

	data.Service().DB.Raw(sqlQuery).Scan(&rawResults)

	seriesTracker := map[int]int{}
	for _, v := range rawResults {
		if _, ok := seriesTracker[v.StationID]; !ok {
			seriesTracker[v.StationID] = len(result.Series)
			result.Series = append(result.Series, v.StationID)
			result.Data = append(result.Data, []int{0, 0, 0, 0, 0, 0, 0})
		}
		dayIdx := extractDayOfWeek(v.StartTime)
		if dayIdx >= 0 && dayIdx < 7 {
			result.Data[seriesTracker[v.StationID]][dayIdx] += v.Secs / 60
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
		ID        int
		StartTime string
		RunSecs   int
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

	cutOff := sqliteDateFilter(days)

	sqlStr := fmt.Sprintf(`SELECT stations.id, station_histories.start_time,
		(duration) as run_secs
	FROM station_histories
	JOIN stations ON stations.id = station_histories.station_id
	WHERE start_time >= %s
		AND stations.id > 0
	ORDER BY station_id ASC`, cutOff)

	data.Service().DB.Raw(sqlStr).Scan(&chartData)
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
		hour := extractHour(i.StartTime)
		if hour >= 0 && hour < 24 {
			result.Data[seriesTracker[i.ID]][hour] += i.RunSecs / 60
		}
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling usage stats.", zap.String("error", err.Error()))
	}

	io.WriteString(rw, string(blob))
}

// statsDailySummary: total watering minutes per day over the range.
func statsDailySummary(rw http.ResponseWriter, req *http.Request) {
	days := parseDaysParam(req, 14)
	type StatsChart struct {
		ReportType int
		Labels     []string
		Series     []string
		Data       [][]int
	}

	type RawResult struct {
		StartTime string
		Mins      int
	}

	var rawResults []RawResult
	cutOff := sqliteDateFilter(days)

	sqlQuery := fmt.Sprintf(`SELECT start_time, SUM(duration)/60 as mins
		FROM station_histories
		WHERE start_time >= %s
		GROUP BY start_time
		ORDER BY start_time ASC`, cutOff)

	data.Service().DB.Raw(sqlQuery).Scan(&rawResults)

	result := StatsChart{
		ReportType: 5,
		Labels:     []string{},
		Series:     []string{"Minutes"},
		Data:       [][]int{[]int{}},
	}

	// ponytail: aggregate by date in Go since strftime() can't parse Go timestamps
	dateMap := map[string]int{}
	var dates []string
	for _, v := range rawResults {
		date := extractDate(v.StartTime)
		if _, exists := dateMap[date]; !exists {
			dates = append(dates, date)
		}
		dateMap[date] += v.Mins
	}
	for _, d := range dates {
		result.Labels = append(result.Labels, d)
		result.Data[0] = append(result.Data[0], dateMap[d])
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		logging.Service().LogError("Error while marshalling daily summary.", zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}
