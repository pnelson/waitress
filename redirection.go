package waitress

import (
	"net/http"
)

func RedirectTo(path string) http.Handler {
	return RedirectToWithCode(path, 303)
}

func RedirectToWithCode(path string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, code)
	})
}
