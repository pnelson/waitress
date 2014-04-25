package router

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// NewConverter is any function that accepts a map of strings to strings and
// returns a Converter.
type NewConverter func(map[string]string) Converter

// A Converter is implemented by objects that can convert their values between
// URL (strings) and proper Go types.
type Converter interface {
	Regexp() string
	Weight() int
	ToGo(value string) (interface{}, error)
	ToUrl(value interface{}) (string, error)
}

// BaseConverter contains common functionality for all converters.
type BaseConverter struct {
	regexp string
	weight int
}

// A StringConverter is the most basic converter as no type conversion is
// necessary. It exists to provide length restrictions.
type StringConverter struct {
	BaseConverter
}

// A PathConverter allows slashes.
type PathConverter struct {
	BaseConverter
}

// An AnyConverter will match against one of the provided options.
type AnyConverter struct {
	BaseConverter
}

// An IntConverter accepts int values. Be careful to not mix this up with the
// Int64Converter which is more common and the default converter for 'int'.
type IntConverter struct {
	BaseConverter
	digits int // The number of fixed digits.
	min    int // The minimum value of the integer.
	max    int // The maximum value of the integer.
}

// An Int64Converter accepts int64 values.
type Int64Converter struct {
	BaseConverter
	digits int64 // The number of fixed digits.
	min    int64 // The minimum value of the integer.
	max    int64 // The maximum value of the integer.
}

// NewStringConverter constructs a new StringConverter from the provided
// arguments. Accepted arguments are: length (exact), minLength, and maxLength.
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

// NewPathConverter constructs a new PathConverter. This converter does not
// accept any arguments.
func NewPathConverter(args map[string]string) Converter {
	if len(args) != 0 {
		return nil
	}
	return &PathConverter{BaseConverter{`[^/].*?`, 200}}
}

// NewAnyConverter constructs a new AnyConverter from the provided 'items'
// argument. 'items' should be a comma-separated string of possible matches.
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

// NewIntConverter constructs a new IntConverter from the provided arguments.
// Accepted arguments are: digits (exact), min, and max. All arguments must be
// of type int.
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

// NewInt64Converter constructs a new Int64Converter from the provided
// arguments.  Accepted arguments are: digits (exact), min, and max. All
// arguments must be of type int64.
func NewInt64Converter(args map[string]string) Converter {
	intArgs, err := mapAtoi64(args)
	if err != nil {
		return nil
	}
	return &Int64Converter{
		BaseConverter{`\d+`, 50},
		intArgs["digits"],
		intArgs["min"],
		intArgs["max"],
	}
}

// Regexp returns the regexp as a string.
func (c *BaseConverter) Regexp() string {
	return c.regexp
}

// Weight returns the converter's weight. The Router uses this to sort.
func (c *BaseConverter) Weight() int {
	return c.weight
}

// ToGo simply returns the provided value as no conversion is necessary.
func (c *BaseConverter) ToGo(value string) (interface{}, error) {
	return value, nil
}

// ToUrl simply returns the provided value as a string.
func (c *BaseConverter) ToUrl(value interface{}) (string, error) {
	return value.(string), nil
}

// ToGo converts the string representation of a number in base 10 to an int.
// If the digits argument was provided during the construction of this
// Int64Converter, then the string length will be checked against. After
// converting, if min or max arguments were provided, the int will be validated
// to be within the provided range.
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

// ToUrl converts the provided value as an int to a string representation in
// base 10. The string will be padded with zero's as necessary based on the
// digits argument used during construction.
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

// ToGo converts the string representation of a number in base 10 to an int64.
// If the digits argument was provided during the construction of this
// Int64Converter, then the string length will be checked against. After
// converting, if min or max arguments were provided, the int64 will be
// validated to be within the provided range.
func (c *Int64Converter) ToGo(value string) (interface{}, error) {
	if c.digits != 0 && int64(len(value)) != c.digits {
		return int64(-1), errors.New("unmatched digits")
	}

	rv, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return int64(-1), err
	}

	if c.min != 0 && rv < c.min || c.max != 0 && rv > c.max {
		return int64(-1), errors.New("not within range")
	}

	return rv, nil
}

// ToUrl converts the provided value as an int64 to a string representation in
// base 10. The string will be padded with zero's as necessary based on the
// digits argument used during construction.
func (c *Int64Converter) ToUrl(value interface{}) (string, error) {
	intValue, ok := value.(int64)
	if !ok {
		return "", errors.New("not a number")
	}

	rv := strconv.FormatInt(intValue, 10)
	if c.digits != 0 {
		rv = fmt.Sprintf(fmt.Sprintf("%%0%ds", c.digits), rv)
	}

	return rv, nil
}

// mapAtoi converts a map of strings representing numbers in base 10 to int.
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

// mapAtoi converts a map of strings representing numbers in base 10 to int64.
func mapAtoi64(args map[string]string) (map[string]int64, error) {
	var err error
	rv := make(map[string]int64)
	for k, v := range args {
		rv[k], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return rv, nil
}
