package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Builder struct {
	handler    http.Handler
	middleware []Middleware
}

func (b *Builder) Use(f Middleware) {
	b.handler = http.Handler(nil)
	b.middleware = append(b.middleware, f)
	for i := len(b.middleware) - 1; i >= 0; i-- {
		b.handler = b.middleware[i](b.handler)
	}
}

func (b *Builder) UseBuilder(builder *Builder) {
	for _, m := range builder.middleware {
		b.Use(m)
	}
}

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

func (b *Builder) UseHandlerFunc(f func(http.ResponseWriter, *http.Request)) {
	b.UseHandler(http.HandlerFunc(f))
}

func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.handler.ServeHTTP(w, r)
}
