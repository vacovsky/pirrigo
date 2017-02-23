package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	//	"time"

	"github.com/davecgh/go-spew/spew"
)

func statsActivityByHour(rw http.ResponseWriter, req *http.Request) {
	type StatsChart struct {
		ReportType int
		Labels     []string
		Series     []string
		Data       [][]float32
	}

	GormDbConnect()
	defer db.Close()

	type RawResult struct {
		Day  int
		Mins float32
	}

	result := StatsChart{
		ReportType: 1,
		Labels:     []string{},
		Series:     []string{},
	}
	if SETTINGS.PirriDebug {
		//		spew.Dump(rawResults1)
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, string(blob))

}

func statsActivityByDayOfWeek(rw http.ResponseWriter, req *http.Request) {
	type StatsChart struct {
		ReportType int
		Labels     []string
		Series     []string
		Data       [][]float32
	}

	result := StatsChart{
		ReportType: 2,
		Labels:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Series:     []string{"Total", "Scheduled", "Unscheduled"},
	}
	GormDbConnect()
	defer db.Close()

	type RawResult struct {
		Day  int
		Mins float32
	}

	var rawResults1 []RawResult
	var rawResults2 []RawResult
	var rawResults3 []RawResult

	sqlQuery1 := fmt.Sprintf(`SELECT DISTINCT DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration / 60) as mins
            FROM station_histories
            WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY)
            GROUP BY day
            ORDER BY day ASC`)
	sqlQuery2 := fmt.Sprintf(`SELECT DISTINCT DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration / 60) as mins
            FROM station_histories
            WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY) AND schedule_id > 0
            GROUP BY day
            ORDER BY day ASC`)
	sqlQuery3 := fmt.Sprintf(`SELECT DISTINCT DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration / 60) as mins
            FROM station_histories
            WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY) AND schedule_id = 0
            GROUP BY day
            ORDER BY day ASC`)
	db.Raw(sqlQuery1, -8, 7).Scan(&rawResults1)
	db.Raw(sqlQuery2, -8, 7).Scan(&rawResults2)
	db.Raw(sqlQuery3, -8, 7).Scan(&rawResults3)

	result.Data = [][]float32{
		[]float32{0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0},
	}

	//	for _, _ := range rawResults1 {
	//				for result.Data[0][result.Data[len(result.Data)-1] != v.Day - 1 {
	//					result.Data[0] = append(result.Data, v.Day - 1)
	//				}

	//	}

	if SETTINGS.PirriDebug {
		//		spew.Dump(rawResults1)
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		fmt.Println(err, err.Error())
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

	GormDbConnect()
	defer db.Close()

	type RawResult struct {
		Day  int
		Mins float32
	}

	result := StatsChart{
		ReportType: 2,
		Labels:     []string{},
		Series:     []string{},
	}
	if SETTINGS.PirriDebug {
		//		spew.Dump(rawResults1)
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, string(blob))
}

func statsStationActivity(rw http.ResponseWriter, req *http.Request) {
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

	sqlStr := fmt.Sprintf(`SELECT stations.id, 
					  HOUR(start_time - INTERVAL %s HOUR) as hour, 
					  (duration) as run_secs
				FROM station_histories
				JOIN stations ON stations.id = station_histories.station_id
				WHERE start_time >= (CURRENT_DATE - INTERVAL 7 DAY) 
				ORDER BY station_id ASC`, strconv.Itoa(-1*SETTINGS.UtcOffset))

	GormDbConnect()
	defer db.Close()

	seriesTracker := map[int]int{}
	db.Raw(sqlStr).Scan(&chartData)

	tracker := 0
	// build list of stations in ascending order.  There must be a better way to do this...
	for y, i := range chartData {
		if tracker == 0 {
			seriesTracker[i.ID] = tracker
			tracker += 1
			result.Data = append(result.Data, []int{
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0})
			result.Series = append(result.Series, i.ID)
		} else if i.ID != result.Series[len(result.Series)-1] {
			seriesTracker[i.ID] = tracker
			result.Data = append(result.Data, []int{
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0})
			result.Series = append(result.Series, i.ID)
			fmt.Println(y)
		}

		result.Data[seriesTracker[i.ID]][i.Hour] += i.RunSecs
	}

	if SETTINGS.PirriDebug {
		spew.Dump(&seriesTracker)
	}
	blob, err := json.Marshal(&result)
	if err != nil {
		fmt.Println(err, err.Error())
	}

	io.WriteString(rw, string(blob))
}
