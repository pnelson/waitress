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
		rule, err := NewRule(path, "", []string{})
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
	_, err := NewRule(`path`, "", []string{})
	if err != ErrLeadingSlash {
		t.Fatalf("unexpected error\nhave %v\nwant %v", err, ErrLeadingSlash)
	}
}

func TestNewRuleMethods(t *testing.T) {
	var methodTests = []struct {
		in  []string
		out []string
	}{
		{[]string{"POST"}, []string{"POST"}},
		{[]string{"POST", "PUT"}, []string{"POST", "PUT"}},
		{[]string{"post"}, []string{"POST"}},
		{[]string{"POST", "POST"}, []string{"POST"}},
		{[]string{"post", "POST"}, []string{"POST"}},
		{[]string{"POST", "post"}, []string{"POST"}},
		{[]string{"GET", "HEAD"}, []string{"GET", "HEAD"}},
		{[]string{"GET"}, []string{"GET", "HEAD"}},
		{[]string{"HEAD"}, []string{"HEAD"}},
	}

	for i, tt := range methodTests {
		rule, err := NewRule(`/`, "", tt.in)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(rule.methods, tt.out) {
			t.Errorf("%d. rule.methods\nhave %v\nwant %v", i, rule.methods, tt.out)
		}
	}
}

func TestRuleDefaults(t *testing.T) {
	args := map[string]interface{}{"foo": "foo", "bar": 4}

	rule, err := NewRule(`/`, "", []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rv := rule.Defaults(args)
	if !reflect.DeepEqual(rule.defaults, args) {
		t.Errorf("rule.Defaults should assign args to `defaults` struct member")
	}
	if rv != rule {
		t.Errorf("rule.Defaults should return itself")
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
		rule, err := NewRule(tt.in, "", []string{})
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

func TestRuleCompileUnbound(t *testing.T) {
	rule, err := NewRule(`/`, "", []string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = rule.compile()
	if err != ErrUnbound {
		t.Errorf("rule.compile\nhave %v\nwant %v", err, ErrUnbound)
	}
}

func TestRuleCompileErrors(t *testing.T) {
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
		rule, err := NewRule(tt.path, "", []string{})
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

func TestRuleMatch(t *testing.T) {
	var matchTests = []struct {
		rule string
		path string
		args args
	}{
		{`/`, "/", args{}},
		{`/<foo>`, "/bar", args{"foo": "bar"}},
		{`/<foo:int>`, "/4", args{"foo": int64(4)}},
		{`/<foo>/<bar>`, "/bar/baz", args{"foo": "bar", "bar": "baz"}},
		{`/<foo>/bar`, "/foo/bar", args{"foo": "foo"}},
	}

	router := New()
	for i, tt := range matchTests {
		rule, err := NewRule(tt.rule, "", []string{})
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

func TestRuleBuild(t *testing.T) {
	var buildTests = []struct {
		rule string
		args args
		out  string
	}{
		{`/`, args{}, "/"},
		{`/foo`, args{}, "/foo"},
		{`/foo`, args{"bar": "bar", "baz": int64(4)}, "/foo?bar=bar&baz=4"},
		{`/foo/<bar>`, args{"bar": "bar", "baz": int64(4)}, "/foo/bar?baz=4"},
		{`/foo/<bar:int(digits=3)>`, args{"bar": int64(4)}, "/foo/004"},
	}

	router := New()
	for i, tt := range buildTests {
		rule, err := NewRule(tt.rule, "", []string{})
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		err = rule.bind(router)
		if err != nil {
			t.Errorf("%d. unexpected error: %v", i, err)
			continue
		}

		out, ok := rule.build(tt.args)
		if !ok || out != tt.out {
			t.Errorf("%d. rule.build(%v)\nhave %q, %t\nwant %q, %t",
				i, tt.args, out, ok, tt.out, true)
			continue
		}
	}
}
