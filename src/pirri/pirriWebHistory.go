package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
	//	"github.com/davecgh/go-spew/spew"
)

func historyAllWeb(rw http.ResponseWriter, req *http.Request) {
	history := []StationHistory{}

	db.Order("id desc").Limit(100).Find(&history)
	blob, err := json.Marshal(&history)
	if err != nil {
		log.LogError("Error while marshalling history from SQL.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"history\": "+string(blob)+"}")
}
