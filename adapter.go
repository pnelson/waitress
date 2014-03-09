package router

import (
	"errors"
)

type Adapter struct {
	router *Router
	method string
	scheme string
	host   string
	path   string
	query  string
}

var ErrNotFound = errors.New("not found")

func NewAdapter(router *Router, method, scheme, host, path, query string) *Adapter {
	return &Adapter{router, method, scheme, host, path, query}
}

func (a *Adapter) Match() (*Rule, map[string]interface{}, error) {
	a.router.sort()

	for _, rule := range a.router.rules {
		args, err := rule.match(a.path)
		if err != nil {
			continue
		}

		return rule, args, nil
	}

	return nil, nil, ErrNotFound
}
