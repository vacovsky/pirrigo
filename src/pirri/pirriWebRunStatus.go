package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vacovsky/pirrigo/src/logging"
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
