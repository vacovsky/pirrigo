package main

func DripNodeById(id int) {

}

type DripNode struct {
	ID    int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPH   float32
	SID   int
	Count int
}

func NewDripnode(gph float32, station int, count int) {
	dn := DripNode{GPH: gph, SID: station, Count: count}
	defer db.Close()
	GormDbConnect()
	db.Create(&dn)
}

func UpdateDripnode(gph float32, station int, count int) {
	dn := DripNode{GPH: gph, SID: station, Count: count}
	defer db.Close()
	GormDbConnect()
	db.Model(&dn).Where(
		"GPH = ?", gph).Where(
		"SID = ?", station).UpdateColumn(DripNode{Count: count})
}
