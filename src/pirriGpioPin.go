package main

import (
	"fmt"
)

type GpioPin struct {
	GPIO  int `gorm:"primary_key"`
	Notes string
}

func GetAllGpio() {
	//	gpio := db.First(&GpioPins)
	fmt.Println("urrrr")
}
