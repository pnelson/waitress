package waitress

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, 200, false}
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.WriteHeader(w.status)
	}

	return w.ResponseWriter.Write(b)
}
