package router

import (
	"testing"
)

func TestNewAdapter(t *testing.T) {
	r := New()
	a := NewAdapter(r, "GET", "http", "localhost", "/", "")
	if a == nil {
		t.Error("router.NewAdapter returned nil")
	}
}
