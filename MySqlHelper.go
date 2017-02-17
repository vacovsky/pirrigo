package main

//import (
//	"database/sql"

//	_ "github.com/go-sql-driver/mysql"
//)

//var db *sql.DB
//var err error

//func openSqlConnection() {
//	db, err = sql.Open(
//		CONFIG.SqlDbType, CONFIG.SqlUser+":"+CONFIG.SqlPass+"@tcp("+CONFIG.SqlServer+":"+CONFIG.SqlPort+")/"+CONFIG.SqlDbName)
//	if err != nil {
//		panic(err.Error())
//	}
//}

//func sqlSmoke() {
//	// sql.DB should be long lived "defer" closes it once this function ends
//	defer db.Close()
//	openSqlConnection()

//	// Test the connection to the database
//	err = db.Ping()
//	if err != nil {
//		panic(err.Error())
//	}
//}

//func sqlSelect(query string) {
//	//	SELECT id, gph, sid, count FROM pirri.dripnodes;
//	//	err := db.QueryRow("SELECT sid FROM  WHERE username=?",
//	//		username).Scan(&databaseUsername, &databasePassword)
//}
