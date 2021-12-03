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
		(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

		// (w).Header().Set("Access-Control-Allow-Headers", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		// (w).Header().Set("Accept", "application/json")
		// (w).Header().Set("User-Agent", "*")
		// (w).Header().Set("Content-Type", "application/json")
		// (w).Header().Set("Referer", "*")

		// Accept: applicatoin/json
		// Content-Type: applicatoin/json
		// Referer: http://localhost:4200/
		// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36

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
