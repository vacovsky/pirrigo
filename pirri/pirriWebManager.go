package pirri

import (
	"net/http"
	"runtime"

	"github.com/vacovsky/pirrigo/settings"
)

func StartPirriWebApp() {
	set := settings.Service()

	for k, v := range protectedRoutes {
		// wrap each route and function in auth handler
		http.HandleFunc(k, basicAuth(v))
	}
	for k, v := range unprotectedRoutes {
		http.HandleFunc(k, v)
	}
	// static content does not require authentication
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// routes to the login page if not authenticated, to the main /home otherwise
	http.HandleFunc("/login", loginAuth)

	// Host server
	panic(http.ListenAndServe(":"+set.Web.Port, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
