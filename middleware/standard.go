package middleware

import (
	"net/http"
	"time"
)

func StripPrefix(prefix string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	}
}

func Timeout(dt time.Duration, msg string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, dt, msg)
	}
}
