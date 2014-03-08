package router

import (
	"fmt"
	"net/http"
	"sort"
)

type Router struct {
	// Map of variable converters available for use in rule paths.
	Converters map[string]NewConverter

	rules  []*Rule // The sequence of rules for this router.
	sorted bool    // Indicates whether or not the rules are already sorted.
}

type sortRules []*Rule

func New() *Router {
	return &Router{
		Converters: map[string]NewConverter{
			"default": NewStringConverter,
			"string":  NewStringConverter,
			"path":    NewPathConverter,
			"any":     NewAnyConverter,
			"int":     NewIntConverter,
		},
	}
}

func (r *Router) Bind(method, scheme, host, path, query string) *Adapter {
	return NewAdapter(r, method, scheme, host, path, query)
}

func (r *Router) BindSimple(scheme, host string) *Adapter {
	return r.Bind("GET", scheme, host, "", "")
}

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

func (r *Router) Rule(path string) (*Rule, error) {
	rule, err := NewRule(path)
	if err != nil {
		return nil, err
	}
	rule.bind(r)
	r.rules = append(r.rules, rule)
	return rule, nil
}

func (r *Router) sort() {
	if r.sorted {
		return
	}
	sort.Sort(sortRules(r.rules))
	r.sorted = true
}

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
