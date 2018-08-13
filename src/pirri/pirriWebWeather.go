package pirri

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vacovsky/pirrigo/src/logging"
	"github.com/vacovsky/pirrigo/src/weather"

	"go.uber.org/zap"
)

func weatherCurrentWeb(rw http.ResponseWriter, req *http.Request) {
	w := weather.Service().Current()
	blob, err := json.Marshal(w)
	if err != nil {
		logging.Service().LogError("Unable to get current weather.",
			zap.String("error", err.Error()))
	}
	io.WriteString(rw, "{ \"weather\": "+string(blob)+"}")
}
