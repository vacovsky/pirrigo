package main

type GpioPin struct {
	GPIO  int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Notes string
}
