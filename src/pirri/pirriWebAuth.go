package main

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

func loginCheck(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Congrats, you're logged in!", 200)
	// getLogger().Debug("Successful login by ")

}

func webHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func loginAuth(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "templates/login.html")
}

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573, modified to also check cookie for auth stuff
// TODO: Clean this sucker up
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		// spew.Dump(r.Header.Get("Authorization"))

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {

			// try cookie auth!
			c, _ := r.Cookie("Authorization")
			// getLogger().LogError("Unable to parse authorization cookie.", zap.String("error", err.Error()))
			q, err := url.ParseQuery(c.Value)
			// getLogger().LogError("Unable to parse query string.", zap.String("error", err.Error()))
			for k := range q {
				s = strings.SplitN(k, " ", 2)
			}
			if len(s) != 2 || err != nil {
				http.Error(w, err.Error(), 401)
				getLogger().LogError("HTTP Authentication Error.", zap.String("error", err.Error()))
				return
			}
		}
		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			getLogger().LogError("HTTP Authentication Error.", zap.String("error", err.Error()))
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if strings.ToLower(pair[0]) != strings.ToLower(SETTINGS.Web.User) || pair[1] != SETTINGS.Web.Secret {
			http.Error(w, "Not authorized", 401)
			getLogger().LogError("HTTP Authentication Error.", zap.String("error", err.Error()))
			return
		}
		h.ServeHTTP(w, r)
	}
}
