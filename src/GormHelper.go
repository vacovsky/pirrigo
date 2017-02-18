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
	db.LogMode(SETTINGS.GormDebug)
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

func JsonifySqlResults(input *gorm.DB) []string {
	var result []string = []string{}
	r, _ := json.Marshal(input.Value)
	result = append(result, string(r))
	fmt.Println(string(r))

	return result

}
