package middleware

import (
	"net/http"
	"time"
)

// StripPrefix uses the next layer of the middleware stack as the http.Handler
// parameter of http.StripPrefix.
func StripPrefix(prefix string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	}
}

// Timeout uses the next layer of the middleware stack as the http.Handler
// parameter of http.TimeoutHandler.
func Timeout(dt time.Duration, msg string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, dt, msg)
	}
}
