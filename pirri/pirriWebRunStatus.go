package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
)

func statusRunWeb(rw http.ResponseWriter, req *http.Request) {
	ORQMutex.Lock()
	status := RUNSTATUS
	ORQMutex.Unlock()
	blob, err := json.Marshal(&status)
	if err != nil {
		logging.Service().LogError("Error while marshalling Run Status from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func statusRunCancel(rw http.ResponseWriter, req *http.Request) {
	ORQMutex.Lock()
	RUNSTATUS.Cancel = true
	status := RUNSTATUS
	ORQMutex.Unlock()
	blob, err := json.Marshal(&status)
	if err != nil {
		logging.Service().LogError("Error while marshalling run status from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func statusRunQueue(rw http.ResponseWriter, req *http.Request) {
	ORQMutex.Lock()
	queue := OfflineRunQueue
	ORQMutex.Unlock()
	blob, err := json.Marshal(&queue)
	if err != nil {
		logging.Service().LogError("Error while marshalling Run Status from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func removeJobFromRunQueue(rw http.ResponseWriter, req *http.Request) {
	type queueId struct {
		QueueIndex int
	}
	var qid queueId
	err := json.NewDecoder(req.Body).Decode(&qid)
	if err != nil {
		logging.Service().LogError("Error while removing job from run queue.",
			zap.String("error", err.Error()))
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}
	ORQMutex.Lock()
	if qid.QueueIndex < 0 || qid.QueueIndex >= len(OfflineRunQueue) {
		ORQMutex.Unlock()
		http.Error(rw, "Invalid queue index", http.StatusBadRequest)
		return
	}
	OfflineRunQueue = removeFromSliceByIndex(OfflineRunQueue, qid.QueueIndex)
	queue := OfflineRunQueue
	ORQMutex.Unlock()
	blob, err := json.Marshal(&queue)
	if err != nil {
		logging.Service().LogError("Error while marshalling Run Status from SQL.",
			zap.String("error", err.Error()))
		http.Error(rw, "Error getting updated queue", http.StatusInternalServerError)
		return
	}
	io.WriteString(rw, string(blob))
}

func removeFromSliceByIndex(sl []*Task, s int) []*Task {
	return append(sl[:s], sl[s+1:]...)
}
