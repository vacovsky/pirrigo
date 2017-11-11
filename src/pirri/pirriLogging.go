package main

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

var instance *Logger
var once sync.Once

type Logger struct {
	lock sync.Mutex
}

func getLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			lock: sync.Mutex{},
		}
	})
	return instance
}

func (l *Logger) LogEvent() {

}

func logToFile(message, stacktrace string) {
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
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.ErrorOutputPaths = []string{SETTINGS.Debug.LogPath}
	cfg.OutputPaths = []string{SETTINGS.Debug.LogPath}

	logger, err := cfg.Build()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	stack := zap.Stack("a stack trace")
	logger.Error("test", stack)
	logger.Debug("logger construction succeeded")

	// Output:
	// {"level":"info","message":"logger construction succeeded","foo":"bar"}
}
