package waitress

import (
	"net/http"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Response: w,
		Request:  r,
	}
}

func (ctx *Context) Abort() http.HandlerFunc {
	return ctx.AbortWithCode(400)
}

func (ctx *Context) AbortWithCode(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: look up code and find appropriate handler
		http.NotFound(w, r)
	}
}

func (ctx *Context) InternalServerError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", 500)
	}
}

func (ctx *Context) Redirect(name string, args map[string]interface{}) http.HandlerFunc {
	path := "" // build
	return ctx.RedirectToWithCode(path, 303)
}

func (ctx *Context) RedirectWithCode(name string, args map[string]interface{}, code int) http.HandlerFunc {
	path := "" // build
	return ctx.RedirectToWithCode(path, code)
}

func (ctx *Context) RedirectTo(path string) http.HandlerFunc {
	return ctx.RedirectToWithCode(path, 303)
}

func (ctx *Context) RedirectToWithCode(path string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, code)
	}
}
