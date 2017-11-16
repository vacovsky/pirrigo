package data

import (
	"sync"
	"time"

	"../logging"
	"../settings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

type ORM struct {
	db   *gorm.DB
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
	d.db, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
	d.db.LogMode(set.Debug.GORM)
	if err != nil {
		log.LogError("Unable to connect to SQL.  Trying again in 15 seconds.",
			zap.String("dbType", set.SQL.DBType),
			zap.String("connectionString", set.SQL.ConnectionString),
			zap.String("error", err.Error()))
		for d.db == nil {
			time.Sleep(time.Duration(15) * time.Second)
			d.db, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
			log.LogError("Unable to connect to SQL on second attempt.  Fatal?  Probably.",
				zap.String("dbType", set.SQL.DBType),
				zap.String("connectionString", set.SQL.ConnectionString),
				zap.String("error", err.Error()))
		}
	}
	err = d.db.DB().Ping()
	if err != nil {
		log.LogError("Ping against SQL database failed.",
			zap.String("error", err.Error()))
	}
}

func (d *ORM) init() {
	d.connect()

	d.db.DB().SetMaxIdleConns(10)
	d.db.DB().SetMaxOpenConns(100)
	d.db.DB().SetConnMaxLifetime(time.Second * 300)
}
