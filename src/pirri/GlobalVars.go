package main

import (
	"fmt"
	"log"
	"sync"
)

var SETTINGS Settings
var SQLCONNSTRING string
var RMQCONNSTRING string
var KILL bool = false
var VERSION string = "0.0.1"
var ERR error
var WG sync.WaitGroup

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
