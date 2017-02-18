package main

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func GormDbConnect() {
	db, err = gorm.Open(SETTINGS.SqlDbType, CONNSTRING)
	if err != nil {
		panic("failed to connect database")
	}

}

func GormSetup() {
	GormDbConnect()
	defer db.Close()
	db.AutoMigrate(
		&DripNode{},
		&GpioPin{},
		&StationHistory{},
		&StationSchedule{},
		&Settings{},
		&Station{})
}

func JsonifyResults(input *gorm.DB) string {
	result, _ := json.Marshal(input)
	fmt.Println(string(result))
	return string(result)
}

//	db.Select(`SELECT id, station, duration FROM schedule
//                        WHERE (startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
//                                AND enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED)))`)

//sqlQuery := fmt.Sprintf(`SELECT id, station, duration FROM schedule
//                        WHERE (startdate <= CAST(replace(date(NOW()), '-', '') AS UNSIGNED)
//                                AND enddate > CAST(replace(date(NOW()), '-', '') AS UNSIGNED))
//                            and %s=1
//                            and starttime=%s`,
//		nowTime.Weekday(),
//		fmt.Sprintf("%02d%02d", nowTime.Hour(), nowTime.Minute()))
