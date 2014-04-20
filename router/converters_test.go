package router

import (
	"testing"
)

type cargs map[string]string

func TestBaseConverter(t *testing.T) {
	var baseConverterTests = []string{
		"foo",
		"bar",
		"foo/bar",
	}

	for i, str := range baseConverterTests {
		c := &BaseConverter{}

		toGoResult, err := c.ToGo(str)
		if err != nil {
			t.Errorf("%d. BaseConverter ToGo(%q) unexpected error: %v", i, str, err)
		}
		if toGoResult != str {
			t.Errorf("%d. BaseConverter ToGo(%q)\nhave %q\nwant %q",
				i, str, toGoResult, str)
		}

		toUrlResult, err := c.ToUrl(str)
		if err != nil {
			t.Errorf("%d. BaseConverter ToUrl(%q) unexpected error: %v", i, str, err)
		}
		if toUrlResult != str {
			t.Errorf("%d. BaseConverter ToUrl(%q)\nhave %q\nwant %q",
				i, str, toUrlResult, str)
		}
	}
}

func TestNewStringConverter(t *testing.T) {
	var stringConverterTests = []struct {
		args   cargs
		regexp string
	}{
		{cargs{}, `[^/]{1,}`},
		{cargs{"minLength": "1"}, `[^/]{1,}`},
		{cargs{"minLength": "2"}, `[^/]{2,}`},
		{cargs{"minLength": "1", "maxLength": "4"}, `[^/]{1,4}`},
		{cargs{"minLength": "2", "maxLength": "4"}, `[^/]{2,4}`},
		{cargs{"minLength": "1", "maxLength": "2", "length": "4"}, `[^/]{4}`},
	}

	for i, tt := range stringConverterTests {
		c := NewStringConverter(tt.args)
		if regexp := c.Regexp(); regexp != tt.regexp {
			t.Errorf("%d. NewStringConverter(%v) regexp\nhave `%s`\nwant `%v`",
				i, tt.args, regexp, tt.regexp)
		}
	}
}

func TestNewPathConverter(t *testing.T) {
	var pathConverterRegexp = `[^/].*?`

	args := cargs{"key": "value"}
	c := NewPathConverter(args)
	if c != nil {
		t.Errorf("NewPathConverter(%v) = %v, want <nil>", args, c)
	}

	args = cargs{}
	c = NewPathConverter(args)
	if regexp := c.Regexp(); regexp != pathConverterRegexp {
		t.Errorf("NewPathConverter regexp\nhave `%v`\nwant `%v`",
			regexp, pathConverterRegexp)
	}
}

func TestNewAnyConverter(t *testing.T) {
	var anyConverterTests = []struct {
		args   cargs
		regexp string
	}{
		{cargs{"items": "a"}, `(?:a)`},
		{cargs{"items": "a,b"}, `(?:a|b)`},
		{cargs{"items": "a,b,c"}, `(?:a|b|c)`},
	}

	for i, tt := range anyConverterTests {
		c := NewAnyConverter(tt.args)
		if regexp := c.Regexp(); regexp != tt.regexp {
			t.Errorf("%d. AnyConverter(%v) regexp\nhave `%s`\nwant `%s`",
				i, tt.args, regexp, tt.regexp)
		}
	}
}

func TestNewAnyConverterNil(t *testing.T) {
	var anyConverterTests = []cargs{
		cargs{},
		cargs{"items": ""},
		cargs{"items": ","},
		cargs{"items": "a,"},
		cargs{"items": ",a"},
	}

	for i, args := range anyConverterTests {
		c := NewAnyConverter(args)
		if c != nil {
			t.Errorf("%d. NewPathConverter(%v) = %v, want <nil>", i, args, c)
		}
	}
}

func TestIntConverter(t *testing.T) {
	var intConverterTests = []struct {
		args        cargs
		regexp      string
		toGoParam   string
		toGoResult  int
		toUrlParam  int
		toUrlResult string
	}{
		{cargs{}, `\d+`, "4", 4, 4, "4"},
		{cargs{"digits": "2"}, `\d+`, "44", 44, 44, "44"},
		{cargs{"digits": "2"}, `\d+`, "04", 4, 4, "04"},
		{cargs{"digits": "2"}, `\d+`, "4", -1, 4, "04"},
		{cargs{"min": "3"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "4"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "5"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"max": "4"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"max": "3"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"min": "3", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "4", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "5", "max": "5"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"min": "3", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "3", "max": "4"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "3", "max": "3"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"digits": "2", "min": "3", "max": "4"}, `\d+`, "04", 4, 4, "04"},
		{cargs{"digits": "2", "min": "3", "max": "4"}, `\d+`, "05", -1, 5, "05"},
	}

	for _, tt := range intConverterTests {
		c := NewIntConverter(tt.args)
		if _, ok := c.(*IntConverter); !ok {
			t.Errorf("NewIntConverter(%v) got <nil>", tt.args)
			continue
		}

		if regexp := c.Regexp(); regexp != tt.regexp {
			t.Errorf("IntConverter regexp\nhave `%v`\nwant `%v`", regexp, tt.regexp)
		}

		toGoResult, err := c.ToGo(tt.toGoParam)
		if err != nil && tt.toGoResult != -1 {
			t.Errorf("IntConverter ToGo(%q) unexpected error: %v",
				tt.toGoParam, err)
		}
		if toGoResult != tt.toGoResult {
			t.Errorf("IntConverter ToGo(%q)\nhave %v\nwant %v",
				tt.toGoParam, toGoResult, tt.toGoResult)
		}

		toUrlResult, err := c.ToUrl(tt.toUrlParam)
		if err != nil {
			t.Errorf("IntConverter ToUrl(%v) unexpected error: %v",
				tt.toUrlParam, err)
		}
		if toUrlResult != tt.toUrlResult {
			t.Errorf("IntConverter ToUrl(%v)\nhave %v\nwant %v",
				tt.toUrlParam, toUrlResult, tt.toUrlResult)
		}
	}
}

func TestInt64Converter(t *testing.T) {
	var int64ConverterTests = []struct {
		args        cargs
		regexp      string
		toGoParam   string
		toGoResult  int64
		toUrlParam  int64
		toUrlResult string
	}{
		{cargs{}, `\d+`, "4", 4, 4, "4"},
		{cargs{"digits": "2"}, `\d+`, "44", 44, 44, "44"},
		{cargs{"digits": "2"}, `\d+`, "04", 4, 4, "04"},
		{cargs{"digits": "2"}, `\d+`, "4", -1, 4, "04"},
		{cargs{"min": "3"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "4"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "5"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"max": "4"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"max": "3"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"min": "3", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "4", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "5", "max": "5"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"min": "3", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "3", "max": "4"}, `\d+`, "4", 4, 4, "4"},
		{cargs{"min": "3", "max": "3"}, `\d+`, "4", -1, 4, "4"},
		{cargs{"digits": "2", "min": "3", "max": "4"}, `\d+`, "04", 4, 4, "04"},
		{cargs{"digits": "2", "min": "3", "max": "4"}, `\d+`, "05", -1, 5, "05"},
	}

	for _, tt := range int64ConverterTests {
		c := NewInt64Converter(tt.args)
		if _, ok := c.(*Int64Converter); !ok {
			t.Errorf("NewInt64Converter(%v) got <nil>", tt.args)
			continue
		}

		if regexp := c.Regexp(); regexp != tt.regexp {
			t.Errorf("Int64Converter regexp\nhave `%v`\nwant `%v`", regexp, tt.regexp)
		}

		toGoResult, err := c.ToGo(tt.toGoParam)
		if err != nil && tt.toGoResult != -1 {
			t.Errorf("Int64Converter ToGo(%q) unexpected error: %v",
				tt.toGoParam, err)
		}
		if toGoResult != tt.toGoResult {
			t.Errorf("Int64Converter ToGo(%q)\nhave %v\nwant %v",
				tt.toGoParam, toGoResult, tt.toGoResult)
		}

		toUrlResult, err := c.ToUrl(tt.toUrlParam)
		if err != nil {
			t.Errorf("Int64Converter ToUrl(%v) unexpected error: %v",
				tt.toUrlParam, err)
		}
		if toUrlResult != tt.toUrlResult {
			t.Errorf("Int64Converter ToUrl(%v)\nhave %v\nwant %v",
				tt.toUrlParam, toUrlResult, tt.toUrlResult)
		}
	}
}
