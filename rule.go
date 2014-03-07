package router

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Rule struct {
	router *Router
	path   string
	regexp *regexp.Regexp
	trace  []trace
	weight int
}

type trace struct {
	param bool
	name  string
}

func NewRule(path string) (*Rule, error) {
	if path == "" || path[0] != '/' {
		return nil, errors.New("rules must begin with a leading slash")
	}
	return &Rule{path: path}, nil
}

func (r *Rule) bind(router *Router) error {
	if r.router != nil {
		return errors.New("rule already bound")
	}
	r.router = router
	return r.compile()
}

func (r *Rule) compile() error {
	var parts []string
	var names []string

	if r.router == nil {
		return errors.New("rule not bound")
	}

	for _, segment := range splitPath(r.path) {
		if segment[0] == '<' {
			name, converter, err := r.parseParam(segment)
			if err != nil {
				return errors.New("malformed path")
			}

			for _, v := range names {
				if v == name {
					return errors.New("duplicate variable name")
				}
			}

			part := fmt.Sprintf(`(?P<%s>%s)`, name, converter.Regexp())
			parts = append(parts, part)

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

// Valid parameters are in the form:
//   <var>
//   <var:converter>
//   <var:converter(arg1=val1,arg2=val2,argx=valx)>
func (r *Rule) parseParam(param string) (string, Converter, error) {
	if len(param) < 3 {
		return "", nil, errors.New("param: too short to be valid")
	}

	if param[0] != '<' || param[len(param)-1] != '>' {
		return "", nil, errors.New("param: must surround with '<' and '>'")
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

	last, arguments := more[len(more)-1], more[:len(more)-1]
	if strings.Contains(more, "(") || last != ')' {
		return "", nil, errors.New("malformed converter")
	}

	args, err := r.parseArguments(arguments)
	if err != nil {
		return "", nil, err
	}

	return name, args, nil
}

func (r *Rule) parseArguments(arguments string) (map[string]string, error) {
	if !strings.Contains(arguments, "=") {
		return nil, errors.New("missing keyword arguments")
	}

	args := make(map[string]string)
	parts := strings.Split(arguments, ",")
	for _, arg := range parts {
		pair := strings.Split(arg, "=")
		if len(pair) != 2 {
			return nil, errors.New("malformed arguments")
		}

		key := pair[0]
		args[key] = pair[1]
	}

	return args, nil
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
