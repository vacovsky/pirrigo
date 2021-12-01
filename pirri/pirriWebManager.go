package pirri

import (
	"net/http"
	"os"
	"runtime"
	"strings"
)

// StartPirriWebApp starts the web server
func StartPirriWebApp() {
	for k, v := range protectedRoutes {
		// wrap each route and function in auth handler
		if strings.ToLower(os.Getenv("PIRRIGO_PASSWORD")) != "" {
			http.HandleFunc(k, basicAuth(v))
		} else {
			http.HandleFunc(k, enableCors(v))
		}
	}

	for k, v := range unprotectedRoutes {
		http.HandleFunc(k, v)
	}
	// static content does not require authentication
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// routes to the login page if not authenticated, to the main /home otherwise
	http.HandleFunc("/login", loginAuth)

	// Host server
	panic(http.ListenAndServe(":"+os.Getenv("PIRRIGO_WEB_PORT"), nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
