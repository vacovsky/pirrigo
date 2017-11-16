package data

import (
	"sync"
	"time"

	"../settings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

type ORM struct {
	db *gorm.DB
}

var instance *ORM
var once sync.Once

//Service returns logging service in a singleton
func Service() *ORM {
	once.Do(func() {
		instance = &ORM{
			lock: sync.Mutex{},
		}
		instance.init()
	})
	return instance
}

func (d *ORM) connect() {
	set := settings.Service()
	var err error
	d.db, err = gorm.Open(set.SQL.DBType, SQLConnString)
	d.db.LogMode(SETTINGS.Debug.GORM)
	if err != nil {
		getLogger().LogError("Unable to connect to SQL.  Trying again in 15 seconds.",
			zap.String("dbType", SETTINGS.SQL.DBType),
			zap.String("connectionString", SQLConnString),
			zap.String("error", err.Error()))
		for d.db == nil {
			time.Sleep(time.Duration(15) * time.Second)
			d.db, err = gorm.Open(SETTINGS.SQL.DBType, SQLConnString)
			getLogger().LogError("Unable to connect to SQL on second attempt.  Fatal?  Probably.",
				zap.String("dbType", SETTINGS.SQL.DBType),
				zap.String("connectionString", SQLConnString),
				zap.String("error", err.Error()))
		}
	}
	err = d.db.DB().Ping()
	if err != nil {
		getLogger().LogError("Ping against SQL database failed.",
			zap.String("error", err.Error()))
	}
}

func (d *ORM) init() {
	d.gormDbConnect()

	d.db.DB().SetMaxIdleConns(10)
	d.db.DB().SetMaxOpenConns(100)
	d.db.DB().SetConnMaxLifetime(time.Second * 300)

	d.db.AutoMigrate(
		&Station{},
		&DripNode{},
		&GpioPin{},
		&StationHistory{},
		&StationSchedule{},
		&Settings{},
	)
}

//TODO: remove this later - it's for testing only.
func (d *ORM) firstRunDBSetup() {
	gpios := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28}
	for pin := range gpios {
		d.db.Create(&GpioPin{
			GPIO:   pin,
			Notes:  "",
			Common: false,
		})
	}
}
