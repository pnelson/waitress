package router

import (
	"reflect"
	"testing"
)

func TestNewRule(t *testing.T) {
	var newRuleTest = []string{
		`/`,
		`/foo`,
		`/foo/bar`,
		`/foo/<bar>`,
		`/foo/<bar>/baz`,
		`/foo/<bar:int>`,
		`/foo/<bar:int(digits=4)>`,
	}

	for i, path := range newRuleTest {
		rule, err := NewRule(path)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		if rule.path != path {
			t.Errorf("%d rule.path have %v, want %v", i, rule.path, path)
		}
	}
}

func TestNewRuleError(t *testing.T) {
	_, err := NewRule(`path`)
	if err != ErrLeadingSlash {
		t.Fatalf("unexpected error\nhave %v\nwant %v", err, ErrLeadingSlash)
	}
}

func TestCompile(t *testing.T) {
	var compileTest = []struct {
		in  string
		out string
	}{
		{`/`, `^/$`},
		{`/foo`, `^/foo$`},
		{`/foo/bar`, `^/foo/bar$`},
		{`/foo/<bar>`, `^/foo/(?P<bar>[^/]{1,})$`},
		{`/foo/<bar>/baz`, `^/foo/(?P<bar>[^/]{1,})/baz$`},
		{`/foo/<bar:int>`, `^/foo/(?P<bar>\d+)$`},
		{`/foo/<bar:int(digits=4)>`, `^/foo/(?P<bar>\d+)$`},
	}

	for i, tt := range compileTest {
		router := New()
		rule, err := NewRule(tt.in)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		err = rule.bind(router)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		if rule.regexp == nil {
			t.Errorf("%d. rule.regexp\nhave %v\nwant %v", i, nil, tt.out)
			continue
		}
		if regexp := rule.regexp.String(); regexp != tt.out {
			t.Errorf("%d. rule.regexp\nhave %v\nwant %v", i, rule.regexp, tt.out)
		}
	}
}

func TestCompileUnbound(t *testing.T) {
	rule, err := NewRule(`/`)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = rule.compile()
	if err != ErrUnbound {
		t.Errorf("rule.compile\nhave %v\nwant %v", err, ErrUnbound)
	}
}

func TestCompileErrors(t *testing.T) {
	var compileTests = []struct {
		path string
		err  error
	}{
		{`/<>`, ErrVariableEmpty},
		{`/<foo`, ErrVariableOpen},
		{`/<foo>/<foo>`, ErrVariableDuplicate},
		{`/<foo:int(>`, ErrConverterOpen},
		{`/<foo:int((>`, ErrConverterOpen},
		{`/<foo:int()>`, nil},
		{`/<foo:int(digits)>`, ErrArguments},
		{`/<foo:int(digits=)>`, ErrArguments},
	}

	router := New()
	for i, tt := range compileTests {
		rule, err := NewRule(tt.path)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}
		err = rule.bind(router)
		if err != tt.err {
			t.Errorf("%d. rule.compile\nhave %v\nwant %v", i, err, tt.err)
		}
	}
}

func TestMatch(t *testing.T) {
	var matchTests = []struct {
		rule string
		path string
		args args
	}{
		{`/`, "/", args{}},
		{`/<foo>`, "/bar", args{"foo": "bar"}},
		{`/<foo:int>`, "/4", args{"foo": 4}},
		{`/<foo>/<bar>`, "/bar/baz", args{"foo": "bar", "bar": "baz"}},
		{`/<foo>/bar`, "/foo/bar", args{"foo": "foo"}},
	}

	router := New()
	for i, tt := range matchTests {
		rule, err := NewRule(tt.rule)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		err = rule.bind(router)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		args, err := rule.match(tt.path)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(args, map[string]interface{}(tt.args)) {
			t.Errorf("%d. rule.match(%q)\nhave %v\nwant %v", i, tt.path, args, tt.args)
		}
	}
}
