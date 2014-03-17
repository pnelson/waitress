package waitress

import (
	"net/http"
)

type Context struct {
	w *ResponseWriter
	r *Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w: &ResponseWriter{w},
		r: &Request{r},
	}
}
