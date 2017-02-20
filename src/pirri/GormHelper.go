package main

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func GormDbConnect() {

	db, ERR = gorm.Open(SETTINGS.SqlDbType, SQLCONNSTRING)
	db.LogMode(SETTINGS.GormDebug)
	if ERR != nil {
		panic(ERR.Error())
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

func JsonifySqlResults(input *gorm.DB) []string {
	var result []string = []string{}
	r, _ := json.Marshal(input.Value)
	result = append(result, string(r))
	fmt.Println(string(r))
	return result
}

//TODO remove this later - it's for testing only.
func createJunkData() {
	GormDbConnect()
	defer db.Close()
	db.Create(&Station{
		GPIO:   5,
		Notes:  "",
		Common: false,
	})
}
