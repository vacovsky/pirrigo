package pirri

/*DripNode Describes a drip emitter */
type DripNode struct {
	ID        int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPH       float32
	StationID int
	Count     int
}
