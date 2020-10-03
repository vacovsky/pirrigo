package pirri

import (
	"io"
	"net/http"
)

func metadataWeb(rw http.ResponseWriter, req *http.Request) {
	result := `{`

	result += `}`

	io.WriteString(rw, result)
}
