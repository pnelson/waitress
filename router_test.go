package router

import (
	"net/http"
	"testing"
)

type args map[string]interface{}

func TestNew(t *testing.T) {
	r := New()
	if r.Converters == nil {
		t.Error("router.New should initialize a set of converters")
	}
}

func TestBind(t *testing.T) {
	r := New()
	adapter := r.Bind("GET", "http", "localhost", "/", "")
	if adapter == nil {
		t.Error("router.Bind returned nil")
	}
}

func TestBindSimple(t *testing.T) {
	r := New()
	adapter := r.BindSimple("http", "localhost")
	if adapter == nil {
		t.Error("router.BindSimple returned nil")
	}
}

func TestBindToRequest(t *testing.T) {
	r := New()
	req, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}

	adapter := r.BindToRequest(req)
	if adapter == nil {
		t.Error("router.BindToRequest returned nil")
	}
}

func TestSort(t *testing.T) {
	var sortTestRules = []struct {
		path  string
		index int
	}{
		{"/", 0},
		{"/<foo>", 4},
		{"/<foo>/<bar:int>", 2},
		{"/<foo>/<bar:path>", 3},
		{"/<foo>/<bar:path>/baz", 1},
	}

	r := New()
	for i, tt := range sortTestRules {
		_, err := r.Rule(tt.path, "", []string{})
		if err != nil {
			t.Fatalf("%d. router.Rule(%q) %v", i, tt.path, err)
		}
	}

	r.sort()
	for i, tt := range sortTestRules {
		if path := r.rules[tt.index].path; path != tt.path {
			t.Errorf("%d. r.rules[%d].path\nhave `%s`\nwant `%s`",
				i, tt.index, path, tt.path)
		}
	}
}
