package router

type Builder struct {
	adapter   *Adapter
	method    string
	name      string
	arguments map[string]interface{}
}

func NewBuilder(adapter *Adapter, method string, name string) *Builder {
	return &Builder{
		adapter:   adapter,
		method:    method,
		name:      name,
		arguments: make(map[string]interface{}),
	}
}

// Path returns the path?query portion of the URL.
func (b *Builder) Path() (string, bool) {
	return b.adapter.build(b)
}

// Full returns the fully qualified URL.
func (b *Builder) Full() (string, bool) {
	return b.adapter.build(b)
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
