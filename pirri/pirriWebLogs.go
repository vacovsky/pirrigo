package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
)

func logsAllWeb(rw http.ResponseWriter, req *http.Request) {
	logs := logging.Service().LoadJournalCtlLogs()
	type LogEntry struct {
		Index int    `json:"index"`
		Msg   string `json:"message"`
	}
	logEntries := make([]LogEntry, 0, len(logs))
	for n, log := range logs {
		if log != "" {
			logEntries = append(logEntries, LogEntry{Index: n, Msg: log})
		}
	}
	blob, err := json.Marshal(map[string][]LogEntry{"logs": logEntries})
	if err != nil {
		logging.Service().LogError("Error marshalling logs", zap.String("error", err.Error()))
		http.Error(rw, "Error marshalling logs", http.StatusInternalServerError)
		return
	}
	io.WriteString(rw, string(blob))
}
