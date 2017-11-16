package pirri

import (
	"fmt"
	"io"
	"net/http"
)

func metadataWeb(rw http.ResponseWriter, req *http.Request) {
	result := `{`

	result += fmt.Sprintf(`"banner": "%s",`, SETTINGS.Pirri.WelcomeMessage)
	result += fmt.Sprintf(`"version": "%s",`, VERSION)
	result += fmt.Sprintf(`"rabbitServer": "%s",`, SETTINGS.RabbitMQ.Server)
	result += fmt.Sprintf(`"sqlServer": "%s",`, SETTINGS.SQL.Server)
	result += fmt.Sprintf(`"newRelicEnabled": %t,`, SETTINGS.NewRelic.Active)
	result += fmt.Sprintf(`"weatherUnits": "%s",`, SETTINGS.Weather.Units)
	result += fmt.Sprintf(`"debug": %t`, SETTINGS.Debug.Pirri)

	result += `}`

	io.WriteString(rw, result)
}
