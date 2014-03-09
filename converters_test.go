package router

import (
	"testing"
)

type cargs map[string]string

func TestStringConverter(t *testing.T) {
	var stringConverterTests = []struct {
		args        cargs
		regexp      string
		toGoParam   string
		toGoResult  string
		toUrlParam  string
		toUrlResult string
	}{
		{cargs{}, `[^/]{1,}`, "test", "test", "test", "test"},
		{cargs{"minLength": "1"}, `[^/]{1,}`, "test", "test", "test", "test"},
		{cargs{"minLength": "1", "maxLength": "4"}, `[^/]{1,4}`, "test", "test", "test", "test"},
		{cargs{"minLength": "1", "maxLength": "2", "length": "4"}, `[^/]{4}`, "test", "test", "test", "test"},
	}

	for _, tt := range stringConverterTests {
		c := NewStringConverter(tt.args)

		if regexp := c.Regexp(); regexp != tt.regexp {
			t.Errorf("StringConverter regexp expected `%v` but got `%v`",
				tt.regexp, regexp)
		}

		toGoResult, err := c.ToGo(tt.toGoParam)
		if err != nil {
			t.Errorf("StringConverter ToGo(%q) unexpected error: %v",
				tt.toGoParam, err)
		}
		if toGoResult != tt.toGoResult {
			t.Errorf("StringConverter ToGo(%q) expected %q but got %q",
				tt.toGoParam, tt.toGoResult, toGoResult)
		}

		toUrlResult, err := c.ToUrl(tt.toUrlParam)
		if err != nil {
			t.Errorf("StringConverter ToUrl(%q) unexpected error: %v",
				tt.toUrlParam, err)
		}
		if toUrlResult != tt.toUrlResult {
			t.Errorf("StringConverter ToUrl(%q) expected %q but got %q",
				tt.toUrlParam, tt.toUrlResult, toUrlResult)
		}
	}
}

func TestPathConverter(t *testing.T) {
	var pathConverterTests = []struct {
		toGoParam   string
		toGoResult  string
		toUrlParam  string
		toUrlResult string
	}{
		{"foo", "foo", "foo", "foo"},
		{"foo/bar", "foo/bar", "foo/bar", "foo/bar"},
	}

	args := cargs{}
	for _, tt := range pathConverterTests {
		c := NewPathConverter(args)

		toGoResult, err := c.ToGo(tt.toGoParam)
		if err != nil {
			t.Errorf("PathConverter ToGo(%q) unexpected error: %v",
				tt.toGoParam, err)
		}
		if toGoResult != tt.toGoResult {
			t.Errorf("PathConverter ToGo(%q) expected %q but got %q",
				tt.toGoParam, tt.toGoResult, toGoResult)
		}

		toUrlResult, err := c.ToUrl(tt.toUrlParam)
		if err != nil {
			t.Errorf("PathConverter ToUrl(%q) unexpected error: %v",
				tt.toUrlParam, err)
		}
		if toUrlResult != tt.toUrlResult {
			t.Errorf("PathConverter ToUrl(%q) expected %q but got %q",
				tt.toUrlParam, tt.toUrlResult, toUrlResult)
		}
	}
}

func TestPathConverterNil(t *testing.T) {
	args := cargs{"key": "value"}
	c := NewPathConverter(args)
	if c != nil {
		t.Errorf("NewPathConverter(%v) = %v, want <nil>", args, c)
	}
}

func TestPathConverterRegexp(t *testing.T) {
	args := cargs{}
	expectedRegexp := `[^/].*?`
	c := NewPathConverter(args)
	if regexp := c.Regexp(); regexp != expectedRegexp {
		t.Errorf("PathConverter regexp expected `%v` but got `%v`",
			expectedRegexp, regexp)
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
			t.Errorf("IntConverter regexp expected `%v` but got `%v`",
				tt.regexp, regexp)
		}

		toGoResult, err := c.ToGo(tt.toGoParam)
		if err != nil && tt.toGoResult != -1 {
			t.Errorf("IntConverter ToGo(%q) unexpected error: %v",
				tt.toGoParam, err)
		}
		if toGoResult != tt.toGoResult {
			t.Errorf("IntConverter ToGo(%q) expected %v but got %v",
				tt.toGoParam, tt.toGoResult, toGoResult)
		}

		toUrlResult, err := c.ToUrl(tt.toUrlParam)
		if err != nil {
			t.Errorf("IntConverter ToUrl(%v) unexpected error: %v",
				tt.toUrlParam, err)
		}
		if toUrlResult != tt.toUrlResult {
			t.Errorf("IntConverter ToUrl(%v) expected %v but got %q",
				tt.toUrlParam, tt.toUrlResult, toUrlResult)
		}
	}
}
