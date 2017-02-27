package main

import (
	"encoding/json"
	"fmt"

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
}

func gormSetup() {
	gormDbConnect()
	defer db.Close()
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
	defer db.Close()
	gormDbConnect()
	db.Create(&Station{
		GPIO:   5,
		Notes:  "",
		Common: false,
	})
}
