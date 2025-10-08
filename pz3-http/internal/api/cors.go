package api

import (
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodPost {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		next.ServeHTTP(w, r)
	})
}
