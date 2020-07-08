package pirri

import (
	"net/http"
	"runtime"

	newrelic "github.com/newrelic/go-agent"
	"github.com/vacovsky/pirrigo/logging"
	"github.com/vacovsky/pirrigo/settings"
)

func StartPirriWebApp() {
	log := logging.Service()
	set := settings.Service()
	if set.NewRelic.Active {
		config := newrelic.NewConfig("PirriGo v"+set.Pirri.Version, set.NewRelic.Key)
		NRAPPMON, err := newrelic.NewApplication(config)

		if NRAPPMON == nil || err != nil {
			log.LogEvent("NewRelic being used.")
		} else {
			for k, v := range protectedRoutes {
				// wrap each route and function in auth handler and new relic
				http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, k, basicAuth(v)))

			}

			for k, v := range unprotectedRoutes {
				// wrap each route and function with new relic
				http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, k, v))

			}
			// static content does not require authentication
			http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, "/static/", func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, r.URL.Path[1:])
			}))

			// routes to the login page if not authenticated, to the main /home otherwise
			http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, "/", loginAuth))
		}
	} else {
		for k, v := range protectedRoutes {
			log.LogEvent("Not using New Relic for: " + k)
			// wrap each route and function in auth handler
			http.HandleFunc(k, basicAuth(v))
		}
		for k, v := range unprotectedRoutes {
			log.LogEvent("Not using New Relic for: " + k)
			http.HandleFunc(k, v)
		}
		// static content does not require authentication
		http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		})

		// routes to the login page if not authenticated, to the main /home otherwise
		http.HandleFunc("/login", loginAuth)
	}

	// Host server
	panic(http.ListenAndServe(":"+set.Web.Port, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
