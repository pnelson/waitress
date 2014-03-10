package router

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type NewConverter func(map[string]string) Converter

type Converter interface {
	Regexp() string
	Weight() int
	ToGo(value string) (interface{}, error)
	ToUrl(value interface{}) (string, error)
}

type BaseConverter struct {
	regexp string
	weight int
}

type StringConverter struct {
	BaseConverter
}

type PathConverter struct {
	BaseConverter
}

type AnyConverter struct {
	BaseConverter
}

type IntConverter struct {
	BaseConverter
	digits int // The number of fixed digits.
	min    int // The minimum value of the integer.
	max    int // The maximum value of the integer.
}

func NewStringConverter(args map[string]string) Converter {
	var minLength string

	regexp := `[^/]`
	if length, ok := args["length"]; ok {
		regexp += fmt.Sprintf(`{%s}`, length)
		return &StringConverter{BaseConverter{regexp, 100}}
	}

	minLength, ok := args["minLength"]
	if !ok {
		minLength = "1"
	}

	if maxLength, ok := args["maxLength"]; ok {
		regexp += fmt.Sprintf(`{%s,%s}`, minLength, maxLength)
	} else {
		regexp += fmt.Sprintf(`{%s,}`, minLength)
	}

	return &StringConverter{BaseConverter{regexp, 100}}
}

func NewPathConverter(args map[string]string) Converter {
	if len(args) != 0 {
		return nil
	}
	return &PathConverter{BaseConverter{`[^/].*?`, 200}}
}

func NewAnyConverter(args map[string]string) Converter {
	arg := strings.Replace(args["items"], " ", "", -1)
	items := strings.Split(arg, ",")

	for i, v := range items {
		if v == "" {
			return nil
		}
		items[i] = regexp.QuoteMeta(v)
	}

	regexp := fmt.Sprintf(`(?:%s)`, strings.Join(items, `|`))
	return &AnyConverter{BaseConverter{regexp, 100}}
}

func NewIntConverter(args map[string]string) Converter {
	intArgs, err := mapAtoi(args)
	if err != nil {
		return nil
	}
	return &IntConverter{
		BaseConverter{`\d+`, 50},
		intArgs["digits"],
		intArgs["min"],
		intArgs["max"],
	}
}

func (c *BaseConverter) Regexp() string {
	return c.regexp
}

func (c *BaseConverter) Weight() int {
	return c.weight
}

func (c *BaseConverter) ToGo(value string) (interface{}, error) {
	return value, nil
}

func (c *BaseConverter) ToUrl(value interface{}) (string, error) {
	return value.(string), nil
}

func (c *IntConverter) ToGo(value string) (interface{}, error) {
	if c.digits != 0 && len(value) != c.digits {
		return -1, errors.New("unmatched digits")
	}

	rv, err := strconv.Atoi(value)
	if err != nil {
		return -1, err
	}

	if c.min != 0 && rv < c.min || c.max != 0 && rv > c.max {
		return -1, errors.New("not within range")
	}

	return rv, nil
}

func (c *IntConverter) ToUrl(value interface{}) (string, error) {
	intValue, ok := value.(int)
	if !ok {
		return "", errors.New("not a number")
	}

	rv := strconv.Itoa(intValue)
	if c.digits != 0 {
		rv = fmt.Sprintf(fmt.Sprintf("%%0%ds", c.digits), rv)
	}

	return rv, nil
}

func mapAtoi(args map[string]string) (map[string]int, error) {
	var err error
	rv := make(map[string]int)
	for k, v := range args {
		rv[k], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return rv, nil
}
