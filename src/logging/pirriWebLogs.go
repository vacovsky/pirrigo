package pirri

import (
	"io"
	"net/http"

	"../logging"
)

func logsAllWeb(rw http.ResponseWriter, req *http.Request) {
	result := `{ "logs": [`
	logs, _ := logging.Service().TailLogs(25)
	for n, log := range logs {
		result += log
		if n < len(logs)-2 {
			result += ","
		}
	}
	result += "]}"
	io.WriteString(rw, result)
}
