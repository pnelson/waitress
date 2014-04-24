package router

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Rule struct {
	router     *Router
	path       string
	name       string
	methods    []string
	defaults   map[string]interface{}
	regexp     *regexp.Regexp
	arguments  []string
	converters map[string]Converter
	trace      []trace
	weight     int
}

type trace struct {
	param bool
	part  string
}

var (
	ErrBound   = errors.New("rule already bound")
	ErrUnbound = errors.New("rule not bound")
)

var (
	ErrLeadingSlash      = errors.New("rules must begin with a leading slash")
	ErrVariableEmpty     = errors.New("variable must have a name")
	ErrVariableOpen      = errors.New("must surround variable with '<' and '>'")
	ErrVariableDuplicate = errors.New("duplicate variable name")
	ErrConverterOpen     = errors.New("must surround converter with '(' and ')'")
	ErrArguments         = errors.New("malformed key/value argument pairs")
)

var (
	ErrMatch = errors.New("path did not match rule")
)

func NewRule(path, name string, methods []string) (*Rule, error) {
	// Ensure that the path begins with a leading slash.
	if path == "" || path[0] != '/' {
		return nil, ErrLeadingSlash
	}

	rule := &Rule{
		path:       path,
		name:       name,
		converters: make(map[string]Converter),
	}

	// Add GET if no methods were provided.
	if len(methods) == 0 {
		methods = append(methods, "GET")
	}

	// Remove duplicate methods and ensure uppercase.
	exist := make(map[string]bool)
	for _, v := range methods {
		method := strings.ToUpper(v)
		if !exist[method] {
			rule.methods = append(rule.methods, method)
			exist[method] = true
		}
	}

	// Add HEAD if not already provided when GET is present.
	if !exist["HEAD"] && exist["GET"] {
		rule.methods = append(rule.methods, "HEAD")
	}

	return rule, nil
}

func (r *Rule) Defaults(args map[string]interface{}) *Rule {
	r.defaults = args
	return r
}

func (r *Rule) Parameters() []string {
	return r.regexp.SubexpNames()[1:]
}

func (r *Rule) allowed(method string) bool {
	for _, m := range r.methods {
		if m == method {
			return true
		}
	}
	return false
}

func (r *Rule) bind(router *Router) error {
	if r.router != nil {
		return ErrBound
	}
	r.router = router
	return r.compile()
}

func (r *Rule) build(args map[string]interface{}) (*url.URL, bool) {
	parts := []string{}
	processed := []string{}
	for _, trace := range r.trace {
		if trace.param {
			part, err := r.converters[trace.part].ToUrl(args[trace.part])
			if err != nil {
				return &url.URL{}, false
			}
			parts = append(parts, part)
			processed = append(processed, trace.part)
		} else {
			parts = append(parts, trace.part)
		}
	}

	for _, key := range processed {
		delete(args, key)
	}

	q := &url.Values{}
	for k, v := range args {
		q.Set(k, fmt.Sprintf("%v", v))
	}

	rv := &url.URL{}
	rv.Path = fmt.Sprintf("/%s", strings.Join(parts, "/"))
	rv.RawQuery = q.Encode()

	return rv, true
}

func (r *Rule) buildable(method string, args map[string]interface{}) bool {
	// Unable to build rule if the method does not match.
	if !r.allowed(method) {
		return false
	}

	// All required values must be present between args and defaults.
	for _, key := range r.arguments {
		if _, ok := r.defaults[key]; !ok {
			if _, ok := args[key]; !ok {
				return false
			}
		}
	}

	// Ensure default values are skipped or equal to args.
	for k, v := range r.defaults {
		if arg, ok := args[k]; ok {
			if arg != v {
				return false
			}
		}
	}

	return true
}

func (r *Rule) compile() error {
	var parts []string
	var names []string

	if r.router == nil {
		return ErrUnbound
	}

	for _, segment := range splitPath(r.path) {
		if segment[0] == '<' {
			name, converter, err := r.parseParam(segment)
			if err != nil {
				return err
			}

			for _, v := range names {
				if v == name {
					return ErrVariableDuplicate
				}
			}

			part := fmt.Sprintf(`(?P<%s>%s)`, name, converter.Regexp())
			parts = append(parts, part)
			names = append(names, name)

			r.arguments = append(r.arguments, name)
			r.converters[name] = converter
			r.trace = append(r.trace, trace{true, name})
			r.weight += converter.Weight()

			continue
		}

		part := regexp.QuoteMeta(segment)
		parts = append(parts, part)

		r.trace = append(r.trace, trace{false, segment})
		r.weight -= len(segment)
	}

	re := fmt.Sprintf(`^/%s$`, strings.Join(parts, "/"))
	r.regexp = regexp.MustCompile(re)

	return nil
}

func (r *Rule) match(path string) (map[string]interface{}, error) {
	var err error
	rv := make(map[string]interface{})

	match := r.regexp.FindStringSubmatch(path)
	if match == nil {
		return nil, ErrMatch
	}

	for i, key := range r.regexp.SubexpNames() {
		if i == 0 || key == "" {
			continue
		}

		rv[key], err = r.converters[key].ToGo(match[i])
		if err != nil {
			return nil, err
		}
	}

	return rv, nil
}

// Valid parameters are in the form:
//   <var>
//   <var:converter>
//   <var:converter(arg1=val1,arg2=val2,argx=valx)>
func (r *Rule) parseParam(param string) (string, Converter, error) {
	if len(param) < 3 {
		return "", nil, ErrVariableEmpty
	}

	if param[0] != '<' || param[len(param)-1] != '>' {
		return "", nil, ErrVariableOpen
	}

	param = param[1 : len(param)-1]
	parts := strings.SplitN(param, ":", 2)

	if len(parts) < 2 {
		parts = append(parts, "default")
	}

	key, args, err := r.parseConverter(parts[1])
	if err != nil {
		return "", nil, err
	}

	converter, ok := r.router.Converters[key]
	if !ok {
		converter = r.router.Converters["default"]
	}

	return parts[0], converter(args), nil
}

func (r *Rule) parseConverter(converter string) (string, map[string]string, error) {
	parts := strings.SplitN(converter, "(", 2)
	if len(parts) == 1 {
		return parts[0], nil, nil
	}

	name := parts[0]
	more := parts[1]

	if more == "" {
		return "", nil, ErrConverterOpen
	}

	last, arguments := more[len(more)-1], more[:len(more)-1]
	if strings.Contains(more, "(") || last != ')' {
		return "", nil, ErrConverterOpen
	}

	args, err := r.parseArguments(arguments)
	if err != nil {
		return "", nil, err
	}

	return name, args, nil
}

func (r *Rule) parseArguments(arguments string) (map[string]string, error) {
	args := make(map[string]string)
	if arguments == "" {
		return args, nil
	}

	if !strings.Contains(arguments, "=") {
		return nil, ErrArguments
	}

	parts := strings.Split(arguments, ",")
	for _, arg := range parts {
		pair := strings.Split(arg, "=")
		if len(pair) != 2 || pair[1] == "" {
			return nil, ErrArguments
		}

		key := pair[0]
		args[key] = pair[1]
	}

	return args, nil
}

func (r *Rule) String() string {
	bound := "unbound"
	if r.router != nil {
		bound = "bound"
	}
	return fmt.Sprintf("<Rule (%s) path:`%s`>", bound, r.path)
}

func splitPath(path string) []string {
	parts := strings.Split(path, "/")
	if parts[0] == "" {
		parts = parts[1:]
	}
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}
