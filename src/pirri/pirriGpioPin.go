package main

//GpioPin - describes a Raspberry Pi GPIO pin
type GpioPin struct {
	ID     int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO   int `gorm:"not null;unique"`
	Notes  string
	Common bool `sql:"DEFAULT:false" gorm:"not null"`
}

func getGpioByPin() {
	gpio := db.Where("GPIO = ?", 4).Find(&GpioPin{}).Order("GPIO DESC")
	jsonifySQLResults(gpio)
}

func getAllGpio() {

}
