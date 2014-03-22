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

type DispatchFunc func(*Rule, map[string]interface{}) (interface{}, error)

var ErrNotFound = errors.New("not found")

func NewAdapter(router *Router, method, scheme, host, path, query string) *Adapter {
	return &Adapter{router, method, scheme, host, path, query}
}

func (a *Adapter) Build(method, name string, args map[string]interface{}) (string, bool) {
	a.router.sort()

	for _, rule := range a.router.names[name] {
		if rule.buildable(method, args) {
			rv, ok := rule.build(args)
			if !ok {
				continue
			}
			return rv, true
		}
	}

	return "", false
}

func (a *Adapter) Dispatch(f DispatchFunc) (interface{}, error) {
	rule, args, err := a.Match()
	if err != nil {
		return nil, err
	}
	return f(rule, args)
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
