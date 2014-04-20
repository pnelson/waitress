package waitress

import (
	"fmt"
	"reflect"
)

type Fragment struct {
	context  reflect.Type
	actions  []func(*state) error
	bindings map[string]interface{}
}

type state struct {
	app    *Application
	prefix string
	name   string
}

func NewFragment(ctx interface{}) *Fragment {
	return &Fragment{
		context:  reflect.TypeOf(ctx),
		bindings: make(map[string]interface{}),
	}
}

func (f *Fragment) Bind(name string, value interface{}) {
	f.bindings[name] = value
}

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

func (f *Fragment) Route(path, name string, methods []string) {
	f.actions = append(f.actions, func(state *state) error {
		rule := state.prefix + path
		endpoint := fmt.Sprintf("%s.%s", state.name, name)
		return state.app.Router.Route(rule, endpoint, f.context, methods)
	})
}
