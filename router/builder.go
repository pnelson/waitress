package router

import (
	"net/url"
)

// A Builder constructs url.URL's using an Adapter.
type Builder struct {
	adapter   *Adapter
	method    string
	name      string
	arguments map[string]interface{}
}

// NewBuilder returns a new Builder.
func NewBuilder(adapter *Adapter, method string, name string) *Builder {
	return &Builder{
		adapter:   adapter,
		method:    method,
		name:      name,
		arguments: make(map[string]interface{}),
	}
}

// Build attempts to return a populated url.URL from the bound Adapter.
func (b *Builder) Build() (*url.URL, bool) {
	rv, ok := b.adapter.build(b)
	if !ok {
		return &url.URL{}, false
	}

	rv.Scheme = b.adapter.scheme
	rv.Host = b.adapter.host

	return rv, true
}

// Get gets the value associated with key.
func (b *Builder) Get(key string) (interface{}, bool) {
	rv, ok := b.arguments[key]
	if !ok {
		return nil, false
	}
	return rv, true
}

// Set sets the key to value. It replaces any existing values.
func (b *Builder) Set(key string, value interface{}) {
	b.arguments[key] = value
}

// Del deletes the argument associated with key.
func (b *Builder) Del(key string) {
	delete(b.arguments, key)
}
