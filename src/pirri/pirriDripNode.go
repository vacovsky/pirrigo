package main

/*DripNode Describes a drip emitter */
type DripNode struct {
	ID    int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPH   float32
	SID   int
	Count int
}

func newDripnode(gph float32, station int, count int) {
	dn := DripNode{GPH: gph, SID: station, Count: count}
	defer db.Close()
	gormDbConnect()
	db.Create(&dn)
}

func updateDripnode(gph float32, station int, count int) {
	dn := DripNode{GPH: gph, SID: station, Count: count}
	defer db.Close()
	gormDbConnect()
	db.Model(&dn).Where(
		"GPH = ?", gph).Where(
		"SID = ?", station).UpdateColumn(DripNode{Count: count})
}
