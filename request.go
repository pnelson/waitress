package waitress

import (
	"net/http"
)

type Request struct {
	*http.Request
}
