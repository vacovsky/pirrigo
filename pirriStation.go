package main

type Station struct {
	ID     int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO   int
	Notes  string
	Common bool `sql:"DEFAULT:false"`
}
