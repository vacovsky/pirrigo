package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var instance *PirriLogger
var once sync.Once

// PirriLogger is the logging thing
type PirriLogger struct {
	lock   sync.Mutex
	logger *zap.Logger
}

//Service returns logging service in a singleton
func Service() *PirriLogger {
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
	cfg.ErrorOutputPaths = []string{os.Getenv("PIRRIGO_LOG_LOCATION")}
	cfg.OutputPaths = []string{os.Getenv("PIRRIGO_LOG_LOCATION")}

	logger, err := cfg.Build()
	l.logger = logger
	if err != nil {
		panic(err)
	}
}

// LogEvent logs events
func (l *PirriLogger) LogEvent(message string, fields ...zapcore.Field) {
	if os.Getenv("PIRRIGO_LOG_LOCATION") != "" {
		fmt.Println("EVENT: ", message)
		defer l.logger.Sync()
		defer l.lock.Unlock()
		l.lock.Lock()
		fields = append(
			fields,
			[]zapcore.Field{
				zap.String("time", time.Now().Format(os.Getenv("PIRRIGO_DATE_FORMAT"))),
			}...,
		)
		l.logger.Debug(
			message,
			fields...,
		)
	}
}

//LogError logs errors
func (l *PirriLogger) LogError(message string, fields ...zapcore.Field) {

	defer l.logger.Sync()
	defer l.lock.Unlock()
	l.lock.Lock()
	fields = append(
		fields,
		[]zapcore.Field{
			zap.String("time", time.Now().Format(os.Getenv("PIRRIGO_DATE_FORMAT"))),
		}...,
	)
	l.logger.Error(
		message,
		fields...,
	)
}

func (l *PirriLogger) TailLogs(lines int) ([]string, error) {
	defer l.lock.Unlock()
	l.lock.Lock()
	cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", lines), os.Getenv("PIRRIGO_LOG_LOCATION"))
	output, err := cmd.Output()
	if err != nil {
		l.LogError("Failed to tail log file.",
			zap.String("error", err.Error()),
			zap.String("logPath", os.Getenv("PIRRIGO_LOG_LOCATION")),
			zap.Int("tailLines", lines),
		)
		return nil, err
	}
	return strings.Split(string(output), "\n"), nil
}
