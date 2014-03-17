package waitress

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
}
