package middleware

import (
	"net/http"
)

func RequestBouncer(f func(*http.Request) bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !f(r) {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
