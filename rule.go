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

var (
	ErrBound   = errors.New("rule already bound")
	ErrUnbound = errors.New("rule not bound")

	ErrLeadingSlash      = errors.New("rules must begin with a leading slash")
	ErrVariableEmpty     = errors.New("variable must have a name")
	ErrVariableOpen      = errors.New("must surround variable with '<' and '>'")
	ErrVariableDuplicate = errors.New("duplicate variable name")
	ErrConverterOpen     = errors.New("must surround converter with '(' and ')'")
	ErrArguments         = errors.New("malformed key/value argument pairs")
)

func NewRule(path string) (*Rule, error) {
	if path == "" || path[0] != '/' {
		return nil, ErrLeadingSlash
	}
	return &Rule{path: path}, nil
}

func (r *Rule) bind(router *Router) error {
	if r.router != nil {
		return ErrBound
	}
	r.router = router
	return r.compile()
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
