package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuilderUse(t *testing.T) {
	rv := ""

	builder := &Builder{}
	builder.Use(testBuilderMiddleware)
	builder.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rv += "a"
			next.ServeHTTP(w, r)
			rv += "b"
		})
	})
	builder.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rv += "c"
			next.ServeHTTP(w, r)
			rv += "d"
		})
	})
	builder.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rv += "X"
		})
	})

	recorder := httptest.NewRecorder()
	builder.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Code != http.StatusTeapot {
		t.Errorf("code have %d want %d", recorder.Code, http.StatusTeapot)
	}

	expected := "acXdb"
	if rv != expected {
		t.Errorf("body have %q want %q", rv, expected)
	}
}

func TestBuilderUseBuilder(t *testing.T) {
	rv := ""

	other := &Builder{}
	other.Use(testBuilderMiddleware)
	other.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "X"
	})

	builder := &Builder{}
	builder.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rv += "a"
			next.ServeHTTP(w, r)
			rv += "b"
		})
	})
	builder.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rv += "c"
			next.ServeHTTP(w, r)
			rv += "d"
		})
	})
	builder.UseBuilder(other)

	recorder := httptest.NewRecorder()
	builder.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Code != http.StatusTeapot {
		t.Errorf("code have %d want %d", recorder.Code, http.StatusTeapot)
	}

	expected := "acXdb"
	if rv != expected {
		t.Errorf("body have %q want %q", rv, expected)
	}
}

func TestBuilderUseHandler(t *testing.T) {
	rv := ""

	builder := &Builder{}
	builder.UseHandler(http.HandlerFunc(testBuilderHandlerFunc))
	builder.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "a"
	}))
	builder.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "X"
	}))

	recorder := httptest.NewRecorder()
	builder.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Code != http.StatusTeapot {
		t.Errorf("code have %d want %d", recorder.Code, http.StatusTeapot)
	}

	expected := "aX"
	if rv != expected {
		t.Errorf("body have %q want %q", rv, expected)
	}
}

func TestBuilderUseHandlerFunc(t *testing.T) {
	rv := ""

	builder := &Builder{}
	builder.UseHandlerFunc(testBuilderHandlerFunc)
	builder.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "a"
	})
	builder.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "X"
	})

	recorder := httptest.NewRecorder()
	builder.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Code != http.StatusTeapot {
		t.Errorf("code have %d want %d", recorder.Code, http.StatusTeapot)
	}

	expected := "aX"
	if rv != expected {
		t.Errorf("have %q want %q", rv, expected)
	}
}

func TestBuilder(t *testing.T) {
	rv := ""

	builder := &Builder{}
	builder.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "a"
	})
	builder.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rv += "b"
			next.ServeHTTP(w, r)
			rv += "c"
		})
	})
	builder.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "X"
	}))
	builder.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rv += "d"
	})

	recorder := httptest.NewRecorder()
	builder.ServeHTTP(recorder, (*http.Request)(nil))

	expected := "abXdc"
	if rv != expected {
		t.Errorf("have %q want %q", rv, expected)
	}
}

func TestBuilderStandardLibrary(t *testing.T) {
	var builderTests = []struct {
		path string
		body string
	}{
		{"/foo/", ""},
		{"/foo/bar", "bar"},
		{"/bar", "404 page not found\n"},
	}

	builder := &Builder{}
	builder.Use(StripPrefix(`/foo/`))
	builder.UseHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.URL.Path)
	})

	for i, tt := range builderTests {
		req, err := http.NewRequest("GET", "http://example.com"+tt.path, nil)
		if err != nil {
			panic(err)
		}

		recorder := httptest.NewRecorder()
		builder.ServeHTTP(recorder, req)

		if body := recorder.Body.String(); body != tt.body {
			t.Errorf("%d. have %q want %q", i, body, tt.body)
		}
	}
}

func testBuilderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		next.ServeHTTP(w, r)
	})
}

func testBuilderHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}
