package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"

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
	if SETTINGS.Debug.LogPath == "" {
		SETTINGS.Debug.LogPath = "pirrigo.log"
	}
	cfg.ErrorOutputPaths = []string{SETTINGS.Debug.LogPath}
	cfg.OutputPaths = []string{SETTINGS.Debug.LogPath}

	logger, err := cfg.Build()
	l.logger = logger
	if err != nil {
		panic(err)
	}
}

func (l *PirriLogger) LogEvent(message string, fields ...zapcore.Field) {
	if SETTINGS.Debug.Pirri {
		fmt.Println("EVENT: ", message)
		defer l.logger.Sync()
		defer l.lock.Unlock()
		l.lock.Lock()
		fields = append(
			fields,
			[]zapcore.Field{
				zap.String("version", VERSION),
				zap.String("time", time.Now().Format(SETTINGS.Pirri.DateFormat)),
			}...,
		)
		l.logger.Debug(
			message,
			fields...,
		)
	}
}
func (l *PirriLogger) LogError(message string, fields ...zapcore.Field) {
	defer l.logger.Sync()
	defer l.lock.Unlock()
	l.lock.Lock()
	fields = append(
		fields,
		[]zapcore.Field{
			zap.String("version", VERSION),
			zap.String("time", time.Now().Format(SETTINGS.Pirri.DateFormat)),
		}...,
	)
	l.logger.Error(
		message,
		fields...,
	)
}

func (l *PirriLogger) tailLogs(lines int) ([]string, error) {
	defer l.lock.Unlock()
	l.lock.Lock()
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), SETTINGS.Debug.LogPath)
	output, err := cmd.Output()
	if err != nil {
		getLogger().LogError("Failed to tail log file.",
			zap.String("error", err.Error()),
			zap.String("logPath", SETTINGS.Debug.LogPath),
			zap.Int("tailLines", lines),
		)
		return nil, err
	}
	return strings.Split(string(output), "\n"), nil
}
