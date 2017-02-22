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

func statsStationActivity(rw http.ResponseWriter, req *http.Request) {
	result := StatsChart{
		ReportType: 4,
		Labels: []string{"00:00", "01:00", "02:00", "03:00",
			"04:00", "05:00", "06:00", "07:00", "08:00",
			"09:00", "10:00", "11:00", "12:00", "13:00",
			"14:00", "15:00", "16:00", "17:00", "18:00",
			"19:00", "20:00", "21:00", "22:00", "23:00"},
		Series: []string{},
	}

	GormDbConnect()
	defer db.Close()

	blob, err := json.Marshal(&result)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"statsData\": "+string(blob)+"}")

}

func statsActivityByDayOfWeek(rw http.ResponseWriter, req *http.Request) {
	result := StatsChart{
		ReportType: 2,
		Labels:     []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		Series:     []string{"Total", "Scheduled", "Unscheduled"},
		//		Data: 		[][]int{
		//		//			[]int{1,2,3,4,5},
		////			[]int{6,2,4,6,2},
		////			[]int{5,2,9,6,1}
		//		}
	}
	GormDbConnect()
	defer db.Close()

	type RawResult struct {
		Day  string
		Mins int
	}

	var rawResults []RawResult
	sqlQuery := fmt.Sprintf(`SELECT DISTINCT DAYOFWEEK((start_time + INTERVAL ? HOUR)) as day, SUM(duration / 60) as mins
            FROM station_histories
            WHERE start_time >= (CURRENT_DATE - INTERVAL ? DAY)
            GROUP BY day
            ORDER BY day ASC`)
	db.Raw(sqlQuery, 8, 7).Scan(&rawResults)
	for _, v := range rawResults {
		dowNum, err := strconv.Atoi(v.Day)
		if err != nil {
			fmt.Println(err, err.Error())
		}
		v.Day = ConvertSqlDayToDOW(dowNum)
	}

	if SETTINGS.PirriDebug {
		spew.Dump(rawResults)
	}

	blob, err := json.Marshal(&result)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"statsData\": "+string(blob)+"}")
}
