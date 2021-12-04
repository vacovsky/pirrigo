package pirri

import (
	"io"
	"net/http"

	"github.com/vacovsky/pirrigo/logging"
)

func logsAllWeb(rw http.ResponseWriter, req *http.Request) {
	result := `{ "logs": [`
	logs, _ := logging.Service().TailLogs(50)
	for n, log := range logs {
		result += log
		if n < len(logs)-2 {
			result += ","
		}
	}
	result += "]}"
	io.WriteString(rw, result)
}
