package main

import (
	"fmt"
	"time"
)

func GetCurrentTasks() StationScheduleEntry {
	defer db.Close()
	result := StationScheduleEntry{}
	openSqlConnection()
	nowTime := time.Now()

	sqlQuery := fmt.Sprintf(`SELECT id, station, duration FROM schedule
                        WHERE (startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
                                AND enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED))
                            and %s=1
                            and starttime=%s`,
		nowTime.Weekday(),
		fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))

	err := db.QueryRow(sqlQuery).Scan(
		&result.ID,
		&result.SID,
		&result.Duration)

	if err != nil {
		fmt.Println(err.Error())
	}
	return result

}

//def get_current_tasks(self):
//        sqlConn = SqlHelper()
//        if str(self.today_cache['day']) + str(self.today_cache['time']) != self.last_datetime:
//            sqlStr = """SELECT id, station, duration from schedule
//                        WHERE (
//                            startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
//                                and enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
//                            )
//                            and {0}=1
//                            and starttime={1}
//        """.format(
//                self.today_cache['day'],
//                self.today_cache['time']
//            )
//            return sqlConn.read(sqlStr)
//        else:
//            pass
