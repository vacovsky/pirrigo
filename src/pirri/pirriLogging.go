package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

var instance *PirriLogger
var once sync.Once

type PirriLogger struct {
	lock   sync.Mutex
	logger *zap.Logger
}

func getLogger() *PirriLogger {
	once.Do(func() {
		instance = &PirriLogger{
			lock: sync.Mutex{},
		}
		instance.init()
	})
	return instance
}

func (l *PirriLogger) init() {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"initialFields": {"application": "PirriGo"},
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	// cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.ErrorOutputPaths = []string{SETTINGS.Debug.LogPath}
	cfg.OutputPaths = []string{SETTINGS.Debug.LogPath}

	logger, err := cfg.Build()
	l.logger = logger
	if err != nil {
		panic(err)
	}
}

func (l *PirriLogger) LogEvent(message string) {
	if SETTINGS.Debug.Pirri {
		fmt.Println("EVENT: ", message)
		defer l.logger.Sync()
		defer l.lock.Unlock()
		l.lock.Lock()
		l.logger.Debug(
			message,
			zap.String("version", VERSION),
			zap.String("time", time.Now().Format(SETTINGS.Pirri.DateFormat)),
		)
	}
}
func (l *PirriLogger) LogError(message, stackTrace string) {
	defer l.logger.Sync()
	defer l.lock.Unlock()
	l.lock.Lock()
	l.logger.Error(
		message,
		zap.String("version", VERSION),
		zap.String("error", stackTrace),
		zap.String("time", time.Now().Format(SETTINGS.Pirri.DateFormat)),
	)
}
