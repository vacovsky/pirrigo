package main

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func gormDbConnect() {
	var err error
	db, err = gorm.Open(SETTINGS.SQL.DBType, SQLConnString)
	db.LogMode(SETTINGS.Debug.GORM)
	if err != nil {
		getLogger().LogError("Unable to connect to SQL.  Trying again in 15 seconds.",
			zap.String("dbType", SETTINGS.SQL.DBType),
			zap.String("connectionString", SQLConnString),
			zap.String("error", err.Error()))
		for db == nil {
			time.Sleep(time.Duration(15) * time.Second)
			db, err = gorm.Open(SETTINGS.SQL.DBType, SQLConnString)
			getLogger().LogError("Unable to connect to SQL on second attempt.  Fatal?  Probably.",
				zap.String("dbType", SETTINGS.SQL.DBType),
				zap.String("connectionString", SQLConnString),
				zap.String("error", err.Error()))
		}
	}
	err = db.DB().Ping()
	if err != nil {
		getLogger().LogError("Ping against SQL database failed.",
			zap.String("error", err.Error()))
	}
}

func gormSetup() {
	gormDbConnect()

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Second * 300)

	db.AutoMigrate(
		&Station{},
		&DripNode{},
		&GpioPin{},
		&StationHistory{},
		&StationSchedule{},
		&Settings{},
	)
}

func jsonifySQLResults(input *gorm.DB) []string {
	var result = []string{}
	r, err := json.Marshal(input.Value)
	if err != nil {
		getLogger().LogError("Problem parsing SQL results.",
			zap.String("error", err.Error()))
	}
	result = append(result, string(r))
	return result
}

//TODO remove this later - it's for testing only.
func firstRunDBSetup() {
	gpios := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28}
	for pin := range gpios {
		db.Create(&GpioPin{
			GPIO:   pin,
			Notes:  "",
			Common: false,
		})
	}
}
