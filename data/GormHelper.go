package data

import (
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	// This is how the documentation indicated to do it.
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/settings"
)

type ORM struct {
	DB   *gorm.DB
	lock sync.Mutex
}

var instance *ORM
var once sync.Once

// Service returns logging service in a singleton
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
	if err != nil {
		log.LogError("Unable to connect to SQL. Trying again in 15 seconds.",
			zap.String("DB.Type", set.SQL.DBType),
			zap.String("connectionString", set.SQL.ConnectionString),
			zap.String("error", err.Error()))
		time.Sleep(time.Duration(15) * time.Second)
		d.DB, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
		if err != nil {
			log.LogError("Fatal: Unable to connect to SQL on retry.",
				zap.String("DB.Type", set.SQL.DBType),
				zap.String("connectionString", set.SQL.ConnectionString),
				zap.String("error", err.Error()))
			panic("Failed to connect to database: " + err.Error())
		}
	}
	d.DB.LogMode(set.Debug.GORM)

	sqlDB := d.DB.DB()
	if err := sqlDB.Ping(); err != nil {
		log.LogError("Ping against SQL database failed.",
			zap.String("error", err.Error()))
	}
}

func (d *ORM) sqliteConnect() {
	log := logging.Service()
	log.LogEvent("Pirrigo initializing with sqlite3 database at " + os.Getenv("PIRRIGO_DB_PATH"))
	var err error

	if os.Getenv("PIRRIGO_DB_PATH") == "" {
		os.Setenv("PIRRIGO_DB_PATH", "pirri.db")
	}
	d.DB, err = gorm.Open("sqlite3", os.Getenv("PIRRIGO_DB_PATH"))

	if err != nil {
		log.LogError(err.Error())
		panic("Failed to connect to sqlite3 database: " + err.Error())
	}
	if os.Getenv("PIRRIGO_DB_LOGMODE") == "" {
		os.Setenv(`PIRRIGO_DB_LOGMODE`, "ON")
	}
	d.DB.LogMode(os.Getenv("PIRRIGO_DB_LOGMODE") == "ON")
}

func (d *ORM) init() {
	if os.Getenv("PIRRIGO_DB_TYPE") != "mysql" {
		d.sqliteConnect()
	} else {
		d.connect()
		d.DB.DB().SetMaxIdleConns(10)
		d.DB.DB().SetMaxOpenConns(100)
		d.DB.DB().SetConnMaxLifetime(time.Second * 300)
	}
}
