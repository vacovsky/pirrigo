package pirri

import (
	"io"
	"net/http"
	"strconv"

	"github.com/vacovsky/pirrigo/logging"
)

func logsAllWeb(rw http.ResponseWriter, req *http.Request) {
	result := `{ "logs": [`
	logs := logging.Service().LoadJournalCtlLogs()
	for n, log := range logs {
		result += strconv.Itoa(n) + " " + log
		// if n < len(logs)-2 {
		// 	result += ","
		// }
	}
	result += "]}"
	io.WriteString(rw, string(result))
}
