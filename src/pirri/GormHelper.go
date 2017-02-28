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
	db, ERR = gorm.Open(SETTINGS.SQLDBType, SQLCONNSTRING)
	db.LogMode(SETTINGS.GormDebug)
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
	var result []string = []string{}
	r, _ := json.Marshal(input.Value)
	result = append(result, string(r))
	fmt.Println(string(r))
	return result
}

//TODO remove this later - it's for testing only.
func createJunkData() {
	db.Create(&Station{
		GPIO:   5,
		Notes:  "",
		Common: false,
	})
}
