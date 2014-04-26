package waitress

import (
	"net/http"
)

// RedirectTo constructs an http.Handler from a path.
func RedirectTo(path string) http.Handler {
	return RedirectToWithCode(path, 303)
}

// RedirectToWithCode constructs an http.Handler from a path and status code.
func RedirectToWithCode(path string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, code)
	})
}
