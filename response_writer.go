package waitress

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, 200}
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.ResponseWriter.WriteHeader(w.status)
	return w.ResponseWriter.Write(b)
}
