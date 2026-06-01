package pirri

import (
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/vacovsky/pirrigo/logging"
	"go.uber.org/zap"
)

// StartPirriWebApp starts the web server
func StartPirriWebApp() {
	defer WG.Done()
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
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		SPAPath := "pirrigo-spa/dist/pirrigo-spa/" + r.URL.Path[1:]
		if r.URL.Path == "/" {
			SPAPath = "pirrigo-spa/dist/pirrigo-spa/index.html"
		}
		http.ServeFile(w, r, SPAPath)
	})

	// routes to the login page if not authenticated, to the main /home otherwise
	http.HandleFunc("/login", loginAuth)

	// Host server
	err := http.ListenAndServe(":"+os.Getenv("PIRRIGO_WEB_PORT"), nil)
	if err != nil {
		logging.Service().LogError("HTTP server failed", zap.String("error", err.Error()))
	}
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
