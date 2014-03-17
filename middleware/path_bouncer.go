package middleware

import (
	"net/http"
)

func PathBouncer(path string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == path {
				http.NotFound(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
