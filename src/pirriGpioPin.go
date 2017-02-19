package main

//	"fmt"

type GpioPin struct {
	ID    int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO  int `gorm:"not null;unique"`
	Notes string
}

func GetGpioByPin() {
	GormDbConnect()
	defer db.Close()
	gpio := db.Where("GPIO = ?", 4).Find(&GpioPin{}).Order("GPIO DESC")
	JsonifySqlResults(gpio)
}

func GetAllGpio() {
	GormDbConnect()
	defer db.Close()

}
