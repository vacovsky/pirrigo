package main

import (
	"sync"
)

var SETTINGS Settings
var SQLCONNSTRING string
var RMQCONNSTRING string
var KILL bool = false
var VERSION string = "0.0.1"
var ERR error
var WG sync.WaitGroup
