package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
)

func statusRunWeb(rw http.ResponseWriter, req *http.Request) {
	blob, err := json.Marshal(&RUNSTATUS)
	if err != nil {
		logging.Service().LogError("Error while marshalling Run Status from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func statusRunCancel(rw http.ResponseWriter, req *http.Request) {
	blob, err := json.Marshal(&RUNSTATUS)
	if err != nil {
		logging.Service().LogError("Error while marshalling run status from SQL.",
			zap.String("error", err.Error()))
	}
	RUNSTATUS.Cancel = true
	io.WriteString(rw, string(blob))
}

func statusRunQueue(rw http.ResponseWriter, req *http.Request) {
	blob, err := json.Marshal(&OfflineRunQueue)
	if err != nil {
		logging.Service().LogError("Error while marshalling Run Status from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, string(blob))
}

func removeJobFromRunQueue(rw http.ResponseWriter, req *http.Request) {
	// if req.Method == "POST" {
	type queueId struct {
		QueueIndex int
	}
	var qid queueId
	err := json.NewDecoder(req.Body).Decode(&qid)
	if err != nil {
		logging.Service().LogError("Error while removing job from run queue.",
			zap.String("error", err.Error()))
	}
	// ORQMutex.Lock()
	OfflineRunQueue = removeFromSliceByIndex(OfflineRunQueue, qid.QueueIndex)
	// ORQMutex.Unlock()
	blob, err := json.Marshal(&OfflineRunQueue)
	if err != nil {
		logging.Service().LogError("Error while marshalling Run Status from SQL.",
			zap.String("error", err.Error()))
		http.Error(rw, "Error getting updatede queue", http.StatusInternalServerError)
		return
	}
	io.WriteString(rw, string(blob))
	// } else {
	// 	http.Error(rw, "Use POST Method", 400)
	// 	return
	// }
}

func removeFromSliceByIndex(sl []*Task, s int) []*Task {
	return append(sl[:s], sl[s+1:]...)
}
