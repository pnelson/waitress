package middleware

import (
	"net/http"
)

func ParseJSON() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") == "application/json" {
				r.ParseForm()
			}
			next.ServeHTTP(w, r)
		})
	}
}
