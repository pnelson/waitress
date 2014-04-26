package waitress

import (
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter. Currently, this is
// only used so that the response status code can be specified before writing
// the header.
type ResponseWriter struct {
	http.ResponseWriter // The http.ResponseWriter to be written to.

	status  int
	written bool
}

// NewResponseWriter returns a new ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, 200, false}
}

// WriteHeader records the response as written and delegates the header writing
// to the embedded ResponseWriter.
func (w *ResponseWriter) WriteHeader(code int) {
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Write will write the bytes to the ResponseWriter.
func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.WriteHeader(w.status)
	}

	return w.ResponseWriter.Write(b)
}
