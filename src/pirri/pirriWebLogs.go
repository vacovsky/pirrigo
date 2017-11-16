package pirri

import (
	"io"
	"net/http"
)

func logsAllWeb(rw http.ResponseWriter, req *http.Request) {
	result := `{ "logs": [`
	logs, _ := log.tailLogs(25)
	for n, log := range logs {
		result += log
		if n < len(logs)-2 {
			result += ","
		}
	}
	result += "]}"
	io.WriteString(rw, result)
}
