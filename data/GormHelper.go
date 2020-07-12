package data

import (
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	// This is how the documentation indicated to do it.
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/settings"
	"go.uber.org/zap"
)

type ORM struct {
	DB   *gorm.DB
	lock sync.Mutex
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
	log := logging.Service()
	set := settings.Service()
	var err error
	d.DB, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
	d.DB.LogMode(set.Debug.GORM)
	if err != nil {
		log.LogError("Unable to connect to SQL.  Trying again in 15 seconds.",
			zap.String("DB.ype", set.SQL.DBType),
			zap.String("connectionString", set.SQL.ConnectionString),
			zap.String("error", err.Error()))
		for d.DB == nil {
			time.Sleep(time.Duration(15) * time.Second)
			d.DB, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
			log.LogError("Unable to connect to SQL on second attempt.  Fatal?  Probably.",
				zap.String("DB.ype", set.SQL.DBType),
				zap.String("connectionString", set.SQL.ConnectionString),
				zap.String("error", err.Error()))
		}
	}
	err = d.DB.DB().Ping()
	if err != nil {
		log.LogError("Ping against SQL database failed.",
			zap.String("error", err.Error()))
	}
}

func (d *ORM) sqliteConnect() {
	log := logging.Service()
	log.LogEvent("Pirrigo initializing with sqlite3 database at " + os.Getenv("PIRRIGO_DB_PATH"))
	var err error
	d.DB, err = gorm.Open("sqlite3", os.Getenv("PIRRIGO_DB_PATH")+".db")

	if err != nil {
		log.LogError(err.Error())
	}
	if os.Getenv("PIRRIGO_DB_LOGMODE") == "" {
		os.Setenv(`PIRRIGO_DB_LOGMODE`, "ON")
	}
	d.DB.LogMode(os.Getenv("PIRRIGO_DB_LOGMODE") == "ON")
}

func (d *ORM) init() {
	if os.Getenv("PIRRIGO_DB_PATH") != "" {
		d.sqliteConnect()
	} else {
		d.connect()
		d.DB.DB().SetMaxIdleConns(10)
		d.DB.DB().SetMaxOpenConns(100)
		d.DB.DB().SetConnMaxLifetime(time.Second * 300)
	}
}
