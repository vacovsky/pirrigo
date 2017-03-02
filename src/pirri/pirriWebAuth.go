package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != SETTINGS.WebUser && pair[1] != SETTINGS.WebPassword {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}
