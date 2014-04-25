package router

import (
	"net/http"
	"strings"
)

// Redirect constructs an http.Handler from a path and status code.
func Redirect(path string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, code)
	})
}

// NotFound returns a 404 Not Found error as an http.Handler.
func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})
}

// MethodNotAllowed constructs an http.Handler from a slice of allowed methods.
func MethodNotAllowed(allowed []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		http.Error(w, "Method Not Allowed", 405)
	})
}

// InternalServerError returns a 500 Internal Server Error as an http.Handler.
func InternalServerError() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", 500)
	})
}
