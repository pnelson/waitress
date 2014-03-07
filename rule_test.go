package router

import (
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
	if err == nil {
		t.Fatalf("incorrect error: %v", err)
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
}

func TestCompileMalformed(t *testing.T) {
}

func TestCompileDuplicateName(t *testing.T) {
}

func TestCompileParamShort(t *testing.T) {
}

func TestCompileParamSurround(t *testing.T) {
}

func TestCompileConverterMalformed(t *testing.T) {
}

func TestCompileArgumentsMissing(t *testing.T) {
}

func TestCompileArgumentsMalformed(t *testing.T) {
}
