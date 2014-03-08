package router

type Adapter struct {
	router *Router
	method string
	scheme string
	host   string
	path   string
	query  string
}

func NewAdapter(router *Router, method, scheme, host, path, query string) *Adapter {
	return &Adapter{router, method, scheme, host, path, query}
}
