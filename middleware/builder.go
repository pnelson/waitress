/*
Package middleware implements a middleware stack inspired by Rack.

Like Rack, each layer is closed over the others. Unlike Rack, this is not done
every request. Stack setup is intended to be done before starting the server.
The outermost layer is cached.

My primary focus here was to stay as close as possible to the http.Handler
interface provided by the Go standard library. You can see some examples in
standard.go that just simply wrap some standard library middleware functions
into middleware that this package can process.
*/
package middleware

import (
	"net/http"
)

// A Middleware is any function that accepts an http.Handler and returns an
// http.Handler. The bottom of the stack will be a nil http.Handler.
type Middleware func(http.Handler) http.Handler

// A Builder builds the middleware stack.
type Builder struct {
	handler    http.Handler
	middleware []Middleware
}

// Use adds a Middleware function to the stack.
func (b *Builder) Use(f Middleware) {
	b.handler = http.Handler(nil)
	b.middleware = append(b.middleware, f)
	for i := len(b.middleware) - 1; i >= 0; i-- {
		b.handler = b.middleware[i](b.handler)
	}
}

// UseBuilder adds an existing middleware stack to this stack.
func (b *Builder) UseBuilder(builder *Builder) {
	for _, m := range builder.middleware {
		b.Use(m)
	}
}

// UseHandler adds an http.Handler to the stack. When the stack is being
// processed, the next http.Handler will be called after the provided handler.
func (b *Builder) UseHandler(f http.Handler) {
	b.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f.ServeHTTP(w, r)
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	})
}

// UseHandlerFunc is a convenience wrapper over UseHandler.
func (b *Builder) UseHandlerFunc(f func(http.ResponseWriter, *http.Request)) {
	b.UseHandler(http.HandlerFunc(f))
}

// ServeHTTP implements the http.Handler interface. It starts off by invoking
// the outermost layer of the stack.
func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.handler.ServeHTTP(w, r)
}
