package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func gormDbConnect() {
	db, ERR = gorm.Open(SETTINGS.SQL.DBType, SQLCONNSTRING)
	db.LogMode(SETTINGS.Debug.GORM)
	if ERR != nil {
		panic(ERR.Error())
	}
	fmt.Println(db.DB().Ping())
}

func gormSetup() {
	gormDbConnect()

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Second * 300)

	db.AutoMigrate(
		&DripNode{},
		&GpioPin{},
		&StationHistory{},
		&StationSchedule{},
		&Settings{},
		&Station{})
}

func jsonifySQLResults(input *gorm.DB) []string {
	var result = []string{}
	r, _ := json.Marshal(input.Value)
	result = append(result, string(r))
	fmt.Println(string(r))
	return result
}

//TODO remove this later - it's for testing only.
func firstRunDBSetup() {
	gpios := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27}
	for pin := range gpios {
		db.Create(&GpioPin{
			GPIO:   pin,
			Notes:  "",
			Common: false,
		})
	}
	db.Create(&Station{
		GPIO:  5,
		Notes: "",
	})
}
