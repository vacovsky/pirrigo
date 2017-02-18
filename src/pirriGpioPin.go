package main

//	"fmt"

type GpioPin struct {
	GPIO  int `gorm:"primary_key"`
	Notes string
}

func GetGpioByPin() {
	GormDbConnect()
	defer db.Close()
	gpio := db.Where("GPIO = ?", 4).Find(&GpioPin{}).Order("GPIO DESC")
	JsonifySqlResults(gpio)
	//	Model(&dn).Where(
	//		"GPH = ?", gph).Where(
	//		"SID = ?", station).UpdateColumn(DripNode{Count: count})
}

func GetAllGpio() {
	GormDbConnect()
	defer db.Close()
	//	gpios := db.Raw("SELECT gpio FROM gpio_pins").Scan()

	//	for gpios.
	//	JsonifyResults(gpios)
	//	Model(&dn).Where(
	//		"GPH = ?", gph).Where(
	//		"SID = ?", station).UpdateColumn(DripNode{Count: count})
}
