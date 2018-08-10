package pirri

import (
	"fmt"
	"io"
	"net/http"

	"github.com/vacovsky/pirrigo/src/settings"
)

func metadataWeb(rw http.ResponseWriter, req *http.Request) {
	set := settings.Service()
	result := `{`

	result += fmt.Sprintf(`"banner": "%s",`, set.Pirri.WelcomeMessage)
	result += fmt.Sprintf(`"version": "%s",`, set.Pirri.Version)
	result += fmt.Sprintf(`"rabbitServer": "%s",`, set.RabbitMQ.Server)
	result += fmt.Sprintf(`"sqlServer": "%s",`, set.SQL.Server)
	result += fmt.Sprintf(`"newRelicEnabled": %t,`, set.NewRelic.Active)
	result += fmt.Sprintf(`"weatherUnits": "%s",`, set.Weather.Units)
	result += fmt.Sprintf(`"debug": %t`, set.Debug.Pirri)

	result += `}`

	io.WriteString(rw, result)
}
