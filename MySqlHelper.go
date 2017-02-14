package main

import (
	//	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func SqlSmoke() {
	//	db, err = sql.Open(SQLDBTYPE, SQLUSER+":"+SQLPASS+"@"+SQLSERVER+"/"+SQLDB)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	// sql.DB should be long lived "defer" closes it once this function ends
	//	defer db.Close()

	//	// Test the connection to the database
	//	err = db.Ping()
	//	if err != nil {
	//		panic(err.Error())
	//	}
}
