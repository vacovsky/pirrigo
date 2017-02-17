package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func GormSetup() {
	fmt.Println(CONNSTRING)
	db, err := gorm.Open(CONFIG.SqlDbType, CONNSTRING)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&StationSchedule{})

	//	db.Select(`SELECT id, station, duration FROM schedule
	//                        WHERE (startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
	//                                AND enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED)))`)
	defer db.Close()

}

//sqlQuery := fmt.Sprintf(`SELECT id, station, duration FROM schedule
//                        WHERE (startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
//                                AND enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED))
//                            and %s=1
//                            and starttime=%s`,
//		nowTime.Weekday(),
//		fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))
