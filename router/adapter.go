package router

import (
	"net/http"
	"net/url"
)

// An Adapter performs URL matching and building on a Router.
type Adapter struct {
	router *Router
	method string
	scheme string
	host   string
	path   string
	query  string
}

// DispatchFunc is a function that will be called when a rule is matched.
type DispatchFunc func(*Rule, map[string]interface{}) interface{}

// NewAdapter returns a new Adapter bound to the provided URL parts.
func NewAdapter(router *Router, method, scheme, host, path, query string) *Adapter {
	return &Adapter{router, method, scheme, host, path, query}
}

// Build returns a new Builder.
func (a *Adapter) Build(method, name string) *Builder {
	return NewBuilder(a, method, name)
}

// Dispatch will attempt to find a rule that matches the Adapter parts. If a
// match is found, the DispatchFunc will be called and the result will be
// returned. Otherwise, the match will return an appropriate http.Handler for
// the error it encountered.
func (a *Adapter) Dispatch(f DispatchFunc) interface{} {
	rule, args, err := a.Match()
	if err != nil {
		return err
	}
	return f(rule, args)
}

// Match attempts to match the Adapter parts to a Rule on the Router. If a path
// match is found but the methods do not match, the MethodNotAllowedHandler,
// bound to the provided methods, will be returned as an error. If no match is
// found at all, the NotFoundHandler will be returned as an error.
func (a *Adapter) Match() (*Rule, map[string]interface{}, http.Handler) {
	a.router.sort()

	var methods []string
	for _, rule := range a.router.rules {
		// Keep trying until we find a match.
		args, err := rule.match(a.path)
		if err != nil {
			continue
		}

		// If the request method is not allowed, keep trying for other matches.
		if !rule.allowed(a.method) {
			i := len(methods)
			methods = append(methods[:i], append(rule.methods, methods[i:]...)...)
			continue
		}

		// A fully matching rule was found.
		return rule, args, nil
	}

	// One or more rules matched but not for the provided method.
	if methods != nil {
		return nil, nil, a.router.MethodNotAllowedHandler(methods)
	}

	// No rule matched the request.
	return nil, nil, a.router.NotFoundHandler()
}

// build constructs a url.URL from the provided builder.
func (a *Adapter) build(builder *Builder) (*url.URL, bool) {
	a.router.sort()

	for _, rule := range a.router.names[builder.name] {
		if rule.buildable(builder.method, builder.arguments) {
			rv, ok := rule.build(builder.arguments)
			if !ok {
				continue
			}

			return rv, true
		}
	}

	return &url.URL{}, false
}
