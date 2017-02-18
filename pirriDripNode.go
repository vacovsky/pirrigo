package main

func DripNodeById(id int) {

}

type DripNode struct {
	ID    int `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	GPH   float32
	SID   int
	Count int
}

func NewDripnode() {
	dn := DripNode{GPH: 6.0, SID: 23, Count: 6}
	GormDbConnect()
	defer db.Close()
	db.Create(&dn)
}
