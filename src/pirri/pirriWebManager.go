package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func stationRunWeb(rw http.ResponseWriter, req *http.Request) {
	var t Task = Task{Station: Station{}, StationSchedule: StationSchedule{}}
	var msr ManualStationRun
	ERR = json.NewDecoder(req.Body).Decode(&msr)

	GormDbConnect()
	defer db.Close()

	db.Where("id = ?", msr.StationID).Find(&t.Station)
	t.StationSchedule = StationSchedule{Duration: msr.Duration}
	if SETTINGS.PirriDebug {
		spew.Dump(t)
		spew.Dump(msr)
	}
	t.Send()
}

func stationAllWeb(rw http.ResponseWriter, req *http.Request) {
	stations := []Station{}

	GormDbConnect()
	defer db.Close()

	db.Limit(100).Find(&stations)
	blob, err := json.Marshal(&stations)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"stations\": "+string(blob)+"}")
}

func stationGetWeb(rw http.ResponseWriter, req *http.Request) {
	var station Station
	stationId, err := strconv.Atoi(req.URL.Query()["stationid"][0])

	GormDbConnect()
	defer db.Close()

	db.Where("id = ?", stationId).Find(&station)
	blob, err := json.Marshal(&station)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, string(blob))
}

func historyAllWeb(rw http.ResponseWriter, req *http.Request) {
	history := []StationHistory{}

	GormDbConnect()
	defer db.Close()

	db.Find(&history)
	blob, err := json.Marshal(&history)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"history\": "+string(blob)+"}")
}

func stationScheduleAllWeb(rw http.ResponseWriter, req *http.Request) {
	stationSchedules := []StationSchedule{}

	GormDbConnect()
	defer db.Close()

	db.Where("end_date > ? AND start_date <= ?", time.Now(), time.Now()).Find(&stationSchedules).Order("ASC")
	blob, err := json.Marshal(&stationSchedules)
	if err != nil {
		fmt.Println(err, err.Error())
	}
	io.WriteString(rw, "{ \"stationSchedules\": "+string(blob)+"}")
}

func stationScheduleEditWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	ERR = json.NewDecoder(req.Body).Decode(&scheduleItem)

	GormDbConnect()
	defer db.Close()

	if db.NewRecord(&scheduleItem) {
		db.Create(&scheduleItem)
	} else {
		db.Save(&scheduleItem)
	}
	if SETTINGS.PirriDebug {
		spew.Dump(scheduleItem)
	}
	stationScheduleAllWeb(rw, req)
}

func stationScheduleDeleteWeb(rw http.ResponseWriter, req *http.Request) {
	var scheduleItem StationSchedule
	ERR = json.NewDecoder(req.Body).Decode(&scheduleItem)

	GormDbConnect()
	defer db.Close()

	db.Delete(&scheduleItem)

	if SETTINGS.PirriDebug {
		spew.Dump(scheduleItem)
	}
	stationScheduleAllWeb(rw, req)
}

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
