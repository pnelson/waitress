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

func TestRouterBind(t *testing.T) {
	r := New()
	adapter := r.Bind("GET", "http", "localhost", "/", "")
	if adapter == nil {
		t.Error("router.Bind returned nil")
	}
}

func TestRouterBindSimple(t *testing.T) {
	r := New()
	adapter := r.BindSimple("http", "localhost")
	if adapter == nil {
		t.Error("router.BindSimple returned nil")
	}
}

func TestRouterBindToRequest(t *testing.T) {
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

func TestRouterRule(t *testing.T) {
	r := New()
	rule, err := r.Rule("/", "", []string{})
	if rule == nil {
		t.Errorf("router.Rule returned nil rule")
	}
	if err != nil {
		t.Errorf("router.Rule returned error: %v", err)
	}
}

func TestRouterRuleError(t *testing.T) {
	r := New()
	rule, err := r.Rule("", "", []string{})
	if rule != nil {
		t.Errorf("router.Rule returned a rule: %v", rule)
	}
	if err == nil {
		t.Errorf("router.Rule returned nil error")
	}
}

func TestRouterMount(t *testing.T) {
	var rules = []struct {
		path string
		name string
	}{
		{"/", "Index"},
		{"/r2/", "R2.Index"},
		{"/r2/<id:int>", "R2.Show"},
	}

	r := New()
	r.Rule("/", "Index", []string{})

	r2 := New()
	r2.Rule("/", "Index", []string{})
	r2.Rule("/<id:int>", "Show", []string{})

	errors := r.Mount("/r2", "R2", r2)
	if errors != nil {
		t.Errorf("router.Mount returned errors: %v", errors)
	}

	for i, rule := range r.rules {
		if rule.path != rules[i].path {
			t.Errorf("%d. rule.path\nhave %q\nwant %q", i, rule.path, rules[i].path)
		}
		if rule.name != rules[i].name {
			t.Errorf("%d. rule.name\nhave %q\nwant %q", i, rule.name, rules[i].name)
		}
	}
}

func TestRouterSort(t *testing.T) {
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
