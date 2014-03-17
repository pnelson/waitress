package waitress

import (
	"net/http"

	"github.com/pnelson/waitress/router"
)

type Router struct {
	*router.Router
}

func NewRouter() *Router {
	return &Router{router.New()}
}

func (r *Router) Route(path, name string, methods []string) {
}

func (r *Router) Mount(prefix string, fragment fragment) {
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}
