package main

//	"fmt"

type GpioPin struct {
	GPIO  int `gorm:"primary_key"`
	Notes string
}

func GetAllGpio() {
	GormDbConnect()
	defer db.Close()
	gpios := db.Where("GPIO > ?", 0).Find(&GpioPin{}).Order("GPIO DESC")

	JsonifyResults(gpios)
	//	Model(&dn).Where(
	//		"GPH = ?", gph).Where(
	//		"SID = ?", station).UpdateColumn(DripNode{Count: count})
}
