package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"../data"
	"../logging"
	"go.uber.org/zap"
	//	"github.com/davecgh/go-spew/spew"
)

func historyAllWeb(rw http.ResponseWriter, req *http.Request) {
	history := []StationHistory{}

	data.Service().DB.Order("id desc").Limit(100).Find(&history)
	blob, err := json.Marshal(&history)
	if err != nil {
		logging.Service().LogError("Error while marshalling history from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"history\": "+string(blob)+"}")
}
