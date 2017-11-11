package main

import (
	"os"
	"sync"

	"github.com/op/go-logging"
)

type logHelper struct {
	Format logging.Formatter
	Logger *logging.Logger
	// BackEnd *logging.LogBackend
	Mutex sync.Mutex
}

// var format = logging.MustStringFormatter(
// 	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
// )

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.

// Password is just an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
func (l *logHelper) createLogFile() {
	if _, err := os.Stat(SETTINGS.Debug.LogPath); os.IsNotExist(err) {
		os.Create(SETTINGS.Debug.LogPath)
	}
}

func (l *logHelper) NewLogHelper() {
	l.createLogFile()
	l.Format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	l.Logger = logging.MustGetLogger("PirriGo")
	l.Mutex = sync.Mutex{}
}

func (l *logHelper) logEvent() {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	writer, err := os.OpenFile(SETTINGS.Debug.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {

	}
	logging.NewLogBackend(writer, "", 0)

}

func testLogging() {
	loggy := logHelper{}
	loggy.NewLogHelper()

	loggy.Logger.Debugf("debug %s")
	// log.Info("info")
	// log.Notice("notice")
	// log.Warning("warning")
	// log.Error("err")
	// log.Critical("crit")

}

// func (l *logHelper) logStuff() {
// 	// For demo purposes, create two backend for os.Stderr.
// 	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
// 	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

// 	// For messages written to backend2 we want to add some additional
// 	// information to the output, including the used log level and the name of
// 	// the function.

// 	// backend2Formatter := logging.NewBackendFormatter(backend2, format)

// 	// Only errors and more severe messages should be sent to backend1
// 	backend1Leveled := logging.AddModuleLevel(backend1)
// 	backend1Leveled.SetLevel(logging.ERROR, "")

// 	// Set the backends to be used.
// 	l.Logger.SetBackend(backend1Leveled, backend2Formatter)

// 	log.Debugf("debug %s")
// 	log.Info("info")
// 	log.Notice("notice")
// 	log.Warning("warning")
// 	log.Error("err")
// 	log.Critical("crit")
// }
