package main

import (
	"fmt"
	"time"
)

func GetCurrentTasks() StationScheduleEntry {
	result := StationScheduleEntry{}

	defer db.Close()
	openSqlConnection()
	nowTime := time.Now()

	sqlQuery := fmt.Sprintf(`SELECT id, station, duration from schedule
                        WHERE (
                            startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
                                and enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
                            )
                            and %s=1
                            and starttime=%s`, nowTime.Weekday(), nowTime.Format())

	//	err := db.QueryRow().Scan(&databaseUsername, &databasePassword

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
