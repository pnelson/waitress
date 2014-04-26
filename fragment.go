package waitress

import (
	"fmt"
	"reflect"
)

// A Fragment is a mountable application that records a series of actions and
// bindings to apply.
type Fragment struct {
	context  reflect.Type
	actions  []func(*state) error
	bindings map[string]interface{}
}

// A state is used for passing contextual information upon registration.
type state struct {
	app    *Application
	prefix string
	name   string
}

// NewFragment returns a new Fragment.
func NewFragment(ctx interface{}) *Fragment {
	return &Fragment{
		context:  reflect.TypeOf(ctx),
		bindings: make(map[string]interface{}),
	}
}

// Bind records a value to bind to the context by a given name.
func (f *Fragment) Bind(name string, value interface{}) {
	f.bindings[name] = value
}

// Registers the Fragment to the Application under a given URL prefix and name.
// All recorded actions and bindings are applied to the Application.
func (f *Fragment) Register(app *Application, prefix, name string) error {
	state := &state{app: app, prefix: prefix, name: name}
	for _, action := range f.actions {
		err := action(state)
		if err != nil {
			return err
		}
	}

	for _, endpoint := range app.Router.endpoints {
		if endpoint.context == f.context {
			for name, value := range f.bindings {
				endpoint.bindings[name] = value
			}
		}
	}

	return nil
}

// Route records the addition of a new route for when it is registered to an
// application.
func (f *Fragment) Route(path, name string, methods []string) {
	f.actions = append(f.actions, func(state *state) error {
		rule := state.prefix + path
		endpoint := fmt.Sprintf("%s.%s", state.name, name)
		return state.app.Router.Route(rule, endpoint, f.context, methods)
	})
}
