package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/b4b4r07/go-pipe"
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

func (l *PirriLogger) LoadJournalCtlLogs() []string {
	defer l.lock.Unlock()
	var b bytes.Buffer

	if err := pipe.Command(&b,
		exec.Command("journalctl", "-xe"),
		exec.Command("grep", "pirrigo"),
	); err != nil {
		log.Fatal(err)
	}
	// pipe.Command(&b,
	// 	exec.Command("journalctl", "-xe"),
	// 	exec.Command("grep", "pirrigo"),
	// )
	io.Copy(os.Stderr, &b)
	l.lock.Lock()

	// return b.String()
	result := strings.Split(b.String(), "\n")
	log.Println("=======================================", result)

	return reverseLogs(result)
}

func reverseLogs(s []string) []string {
	i := 0
	j := len(s) - 1
	for i < j {
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
	return s
}
