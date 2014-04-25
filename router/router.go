/*
Package router implements a flexible HTTP routing package.

This package was originally inspired by the Werkzeug router. It certainly lacks
a lot of the features and functionality but does provide some nice abstractions
over the DefaultServeMux provided by the Go standard library.

The biggest differentiator between this package and any other Go routing
pacakges that I have seen are the Converters. Most commonly, this is used to
convert integer identification numbers in the URL into an int or int64. This is
useful on its own but the power lies in the flexibility this provides. This
package ships with a number of common converters and you can always build your
own custom converters.
*/
package router

import (
	"fmt"
	"net/http"
	"sort"
)

// A Router stores all of the rules and configuration.
type Router struct {
	// Map of variable converters available for use in rule paths.
	Converters map[string]NewConverter

	// The handler to call when a path must be redirected.
	RedirectHandler func(string, int) http.Handler
	// The handler to call when a path is not matched.
	NotFoundHandler func() http.Handler
	// The handler to call when a path is matched but the HTTP method is not.
	MethodNotAllowedHandler func([]string) http.Handler
	// The handler to call when all hell when something terrible happens.
	InternalServerErrorHandler func() http.Handler

	rules  []*Rule // The sequence of rules for this router.
	sorted bool    // Indicates whether or not the rules are already sorted.

	names map[string][]*Rule // Map of rules by name.
}

type sortRules []*Rule // A thin wrapper used to implement sort.Interface.
type sortNames []*Rule // A thin wrapper used to implement sort.Interface.

// New returns a new Router preconfigured with some common URL param converters
// and sane defaults for HTTP 3xx, 404, 405, and 500 errors.
func New() *Router {
	return &Router{
		Converters: map[string]NewConverter{
			"default": NewStringConverter,
			"string":  NewStringConverter,
			"path":    NewPathConverter,
			"any":     NewAnyConverter,
			"int":     NewInt64Converter,
		},

		RedirectHandler:            Redirect,
		NotFoundHandler:            NotFound,
		MethodNotAllowedHandler:    MethodNotAllowed,
		InternalServerErrorHandler: InternalServerError,

		names: make(map[string][]*Rule),
	}
}

// Bind returns a new Adapter bound to the provided URL parts.
func (r *Router) Bind(method, scheme, host, path, query string) *Adapter {
	return NewAdapter(r, method, scheme, host, path, query)
}

// BindSimple returns a new Adapter with the minimal required parameters.
func (r *Router) BindSimple(scheme, host string) *Adapter {
	return r.Bind("GET", scheme, host, "", "")
}

// BindToRequest returns a new Adapter bound to the provided request.
func (r *Router) BindToRequest(req *http.Request) *Adapter {
	method := req.Method
	scheme := "https"
	if req.TLS == nil {
		scheme = "http"
	}
	host := req.Host
	path := req.URL.Path
	query := req.URL.RawQuery
	return r.Bind(method, scheme, host, path, query)
}

// Rule registers a new rule bound to this router.
func (r *Router) Rule(path, name string, methods []string) (*Rule, error) {
	rule, err := NewRule(path, name, methods)
	if err != nil {
		return nil, err
	}
	rule.bind(r)
	r.rules = append(r.rules, rule)
	r.names[name] = append(r.names[name], rule)
	return rule, nil
}

// Mount registers another router bound to this router under some prefix.
func (r *Router) Mount(prefix, name string, router *Router) []error {
	var errors []error
	for _, rule := range router.rules {
		_, err := r.Rule(
			prefix+rule.path,
			fmt.Sprintf("%s.%s", name, rule.name),
			rule.methods,
		)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// sort will sort the rules if needed.
func (r *Router) sort() {
	if r.sorted {
		return
	}

	sort.Sort(sortRules(r.rules))
	for _, rules := range r.names {
		sort.Sort(sortNames(rules))
	}

	r.sorted = true
}

// String is implemented for debugging purposes and will print the rule map.
func (r *Router) String() string {
	rv := "\n"
	for _, rule := range r.rules {
		rv += fmt.Sprintf("  %s\n", rule)
	}
	return fmt.Sprintf("<Router rules:[%s]>", rv)
}

func (s sortRules) Len() int      { return len(s) }
func (s sortRules) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortRules) Less(i, j int) bool {
	// Rules without arguments come first for performance.
	if len(s[i].arguments) == 0 {
		return true
	}
	if len(s[j].arguments) == 0 {
		return false
	}

	// Rules that are more complex come next.
	if len(s[i].trace) > len(s[j].trace) {
		return true
	}
	if len(s[i].trace) < len(s[j].trace) {
		return false
	}

	// Lastly, rules are sorted by ascending weight.
	return s[i].weight < s[j].weight
}

func (s sortNames) Len() int      { return len(s) }
func (s sortNames) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortNames) Less(i, j int) bool {
	// Rules with more arguments come first.
	if len(s[i].arguments) > len(s[j].arguments) {
		return true
	}
	if len(s[i].arguments) < len(s[j].arguments) {
		return false
	}

	// Lastly, rules are sorted by descending default argument quantity.
	return len(s[i].defaults) > len(s[j].defaults)
}
