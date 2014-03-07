package router

import (
	"testing"
)

type args map[string]string

func TestStringConverter(t *testing.T) {
	var stringConverterTests = []struct {
		args        args
		regexp      string
		toGoParam   string
		toGoResult  string
		toUrlParam  string
		toUrlResult string
	}{
		{args{}, `[^/]{1,}`, "test", "test", "test", "test"},
		{args{"minLength": "1"}, `[^/]{1,}`, "test", "test", "test", "test"},
		{args{"minLength": "1", "maxLength": "4"}, `[^/]{1,4}`, "test", "test", "test", "test"},
		{args{"minLength": "1", "maxLength": "2", "length": "4"}, `[^/]{4}`, "test", "test", "test", "test"},
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

	args := args{}
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
	args := args{"key": "value"}
	c := NewPathConverter(args)
	if c != nil {
		t.Errorf("NewPathConverter(%v) = %v, want <nil>", args, c)
	}
}

func TestPathConverterRegexp(t *testing.T) {
	args := args{}
	expectedRegexp := `[^/].*?`
	c := NewPathConverter(args)
	if regexp := c.Regexp(); regexp != expectedRegexp {
		t.Errorf("PathConverter regexp expected `%v` but got `%v`",
			expectedRegexp, regexp)
	}
}

func TestIntConverter(t *testing.T) {
	var intConverterTests = []struct {
		args        args
		regexp      string
		toGoParam   string
		toGoResult  int
		toUrlParam  int
		toUrlResult string
	}{
		{args{}, `\d+`, "4", 4, 4, "4"},
		{args{"digits": "2"}, `\d+`, "44", 44, 44, "44"},
		{args{"digits": "2"}, `\d+`, "04", 4, 4, "04"},
		{args{"digits": "2"}, `\d+`, "4", -1, 4, "04"},
		{args{"min": "3"}, `\d+`, "4", 4, 4, "4"},
		{args{"min": "4"}, `\d+`, "4", 4, 4, "4"},
		{args{"min": "5"}, `\d+`, "4", -1, 4, "4"},
		{args{"max": "5"}, `\d+`, "4", 4, 4, "4"},
		{args{"max": "4"}, `\d+`, "4", 4, 4, "4"},
		{args{"max": "3"}, `\d+`, "4", -1, 4, "4"},
		{args{"min": "3", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{args{"min": "4", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{args{"min": "5", "max": "5"}, `\d+`, "4", -1, 4, "4"},
		{args{"min": "3", "max": "5"}, `\d+`, "4", 4, 4, "4"},
		{args{"min": "3", "max": "4"}, `\d+`, "4", 4, 4, "4"},
		{args{"min": "3", "max": "3"}, `\d+`, "4", -1, 4, "4"},
		{args{"digits": "2", "min": "3", "max": "4"}, `\d+`, "04", 4, 4, "04"},
		{args{"digits": "2", "min": "3", "max": "4"}, `\d+`, "05", -1, 5, "05"},
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
