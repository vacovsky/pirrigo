package pirri

import "time"

//Station describes a "Zone" or "Station" as used in garden irrigation.
type Station struct {
	ID          int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPIO        int `gorm:"not null;unique"`
	Notes       string
	Name        string
	Description string
	LastRun     time.Time
	NextRun     time.Time
}
