package main

//import (
//	"encoding/json"
//	"fmt"
//	"log"
//	"net/http"
//)

//func test(rw http.ResponseWriter, req *http.Request) {
//	decoder := json.NewDecoder(req.Body)
//	var t test_struct
//	err := decoder.Decode(&t)
//	if err != nil {
//		panic(err)
//	}
//	defer req.Body.Close()
//	fmt.Println("test struct", t)
//	log.Println(t.Test)
//}
