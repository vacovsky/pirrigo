package data

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"../settings"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
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

	// log := logging.Service()
	set := settings.Service()
	var err error
	// if set.SQL.UseSSL {
	// SetSSL(set.SQL.KeyPath, set.SQL.CAPath, set.SQL.CertPath)
	// tlsProp := set.SQL.ConnectionString + "&tls=enabled"
	// log.Println(tlsProp)
	// d.DB, err = gorm.Open(set.SQL.DBType, tlsProp)
	// } else {
	d.DB, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
	// }
	d.DB.LogMode(set.Debug.GORM)
	if err != nil {
		log.Println("Unable to connect to database: ", err)
		// log.LogError("Unable to connect to SQL.  Trying again in 15 seconds.",
		// 	zap.String("DB.ype", set.SQL.DBType),
		// 	zap.String("connectionString", set.SQL.ConnectionString),
		// 	zap.String("error", err.Error()))
		for d.DB == nil {
			time.Sleep(time.Duration(15) * time.Second)
			d.DB, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
			log.Println("Unable to connect to database a second time...  Check your configuration", err)

			// log.LogError("Unable to connect to SQL on second attempt.  Fatal?  Probably.",
			// 	zap.String("DB.type", set.SQL.DBType),
			// 	zap.String("connectionString", set.SQL.ConnectionString),
			// 	zap.String("error", err.Error()))
		}
	}
	err = d.DB.DB().Ping()
	if err != nil {
		// log.LogError("Ping against SQL database failed.",
		// 	zap.String("error", err.Error()))
	}

	// log := logging.Service()
	// set := settings.Service()
	// var err error
	// d.DB, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
	// d.DB.LogMode(set.Debug.GORM)
	// if err != nil {
	// log.LogError("Unable to connect to SQL.  Trying again in 15 seconds.",
	// zap.String("DB.ype", set.SQL.DBType),
	// zap.String("connectionString", set.SQL.ConnectionString),
	// zap.String("error", err.Error()))
	// for d.DB == nil {
	// time.Sleep(time.Duration(15) * time.Second)
	// d.DB, err = gorm.Open(set.SQL.DBType, set.SQL.ConnectionString)
	// log.LogError("Unable to connect to SQL on second attempt.  Fatal?  Probably.",
	// zap.String("DB.ype", set.SQL.DBType),
	// zap.String("connectionString", set.SQL.ConnectionString),
	// zap.String("error", err.Error()))
	// }
	// }
	// err = d.DB.DB().Ping()
	// if err != nil {
	// log.LogError("Ping against SQL database failed.",
	// zap.String("error", err.Error()))
	// }
}

func SetSSL(keyPath, caPath, certPath string) {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(caPath)
	if err != nil {
		log.Fatal(err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}
	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatal(err)
	}
	clientCert = append(clientCert, certs)
	mysql.RegisterTLSConfig("enabled", &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	})
}
func (d *ORM) init() {
	d.connect()

	d.DB.DB().SetMaxIdleConns(10)
	d.DB.DB().SetMaxOpenConns(100)
	d.DB.DB().SetConnMaxLifetime(time.Second * 300)
}
