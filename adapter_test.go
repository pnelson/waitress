package router

import (
	"reflect"
	"testing"
)

func TestNewAdapter(t *testing.T) {
	r := New()
	a := NewAdapter(r, "GET", "http", "localhost", "/", "")
	if a == nil {
		t.Error("router.NewAdapter returned nil")
	}
}

func TestAdapterMatchRule(t *testing.T) {
	var matchTests = []struct {
		// in
		method string
		host   string
		path   string

		// out
		want string
	}{
		{"GET", "localhost", "/", "/"},
		{"GET", "localhost", "/a", "/a"},
		{"GET", "localhost", "/a/foo", "/a/<foo>"},
		{"GET", "localhost", "/a/foo/bar", "/a/<foo>/<bar:path>"},
		{"GET", "localhost", "/a/foo/bar/qux", "/a/<foo>/<bar:path>"},
		{"GET", "localhost", "/a/foo/bar/baz", "/a/<foo>/<bar:path>/baz"},
		{"GET", "localhost", "/a/4", "/a/<foo:int>"},
	}

	r := basicAdapterSetup(t)
	for i, tt := range matchTests {
		adapter := r.Bind(tt.method, "http", tt.host, tt.path, "")
		if adapter == nil {
			t.Errorf("%d. problem binding with %s http://%s%s",
				i, tt.method, tt.host, tt.path)
			continue
		}

		rule, _, err := adapter.Match()
		if err != nil {
			t.Errorf("%d. unexpected error %v", i, err)
			continue
		}

		if rule.path != tt.want {
			t.Errorf("%d. adapter.Match\nhave `%s`\nwant `%s`", i, rule.path, tt.want)
		}
	}
}

func TestAdapterMatchArgs(t *testing.T) {
	var matchTests = []struct {
		// in
		method string
		host   string
		path   string

		// out
		want args
	}{
		{"GET", "localhost", "/", args{}},
		{"GET", "localhost", "/a", args{}},
		{"GET", "localhost", "/a/foo", args{"foo": "foo"}},
		{"GET", "localhost", "/a/foo/bar", args{"foo": "foo", "bar": "bar"}},
		{"GET", "localhost", "/a/foo/bar/qux", args{"foo": "foo", "bar": "bar/qux"}},
		{"GET", "localhost", "/a/foo/bar/baz", args{"foo": "foo", "bar": "bar"}},
		{"GET", "localhost", "/a/4", args{"foo": 4}},
	}

	r := basicAdapterSetup(t)
	for i, tt := range matchTests {
		adapter := r.Bind(tt.method, "http", tt.host, tt.path, "")
		if adapter == nil {
			t.Errorf("%d. problem binding with %s http://%s%s",
				i, tt.method, tt.host, tt.path)
			continue
		}

		_, args, err := adapter.Match()
		if err != nil {
			t.Errorf("%d. unexpected error %v", i, err)
			continue
		}

		if !reflect.DeepEqual(args, map[string]interface{}(tt.want)) {
			t.Errorf("%d. adapter.Match\nhave %v\nwant %v", i, args, tt.want)
		}
	}
}

func TestAdapterMatchErrors(t *testing.T) {
	var matchTests = []struct {
		// in
		method string
		host   string
		path   string

		// out
		err error
	}{
		{"GET", "localhost", "/b", ErrNotFound},
	}

	r := basicAdapterSetup(t)
	for i, tt := range matchTests {
		adapter := r.Bind(tt.method, "http", tt.host, tt.path, "")
		if adapter == nil {
			t.Errorf("%d. problem binding with %s http://%s%s",
				i, tt.method, tt.host, tt.path)
			continue
		}

		_, _, err := adapter.Match()
		if err != tt.err {
			t.Errorf("%d. adapter.Match\nhave %v\nwant %v", i, err, tt.err)
		}
	}
}

func basicAdapterSetup(t *testing.T) *Router {
	var basicRules = []string{
		"/",
		"/a",
		"/a/<foo>",
		"/a/<foo>/<bar:path>",
		"/a/<foo>/<bar:path>/baz",
		"/a/<foo:int>",
	}

	r := New()
	for i, path := range basicRules {
		_, err := r.Rule(path, []string{})
		if err != nil {
			t.Fatalf("%d. router.Rule(%q) %v", i, path, err)
		}
	}

	return r
}
