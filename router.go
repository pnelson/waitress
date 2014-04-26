package waitress

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/pnelson/waitress/router"
)

// A Router is a wrapper over the waitress/router package.
type Router struct {
	*router.Router // The waitress/router Router is embedded for its methods.

	endpoints map[*router.Rule]*endpoint
}

// An endpoint needs to keep track of the context it belongs to, the method it
// will be calling, and any additional bindings to apply to the context.
type endpoint struct {
	context  reflect.Type
	method   reflect.Value
	bindings map[string]interface{}
}

// NewRouter returns a new Router.
func NewRouter() *Router {
	return &Router{
		Router:    router.New(),
		endpoints: make(map[*router.Rule]*endpoint),
	}
}

// Route registers a route by method name.
func (r *Router) Route(path, name string, context reflect.Type, methods []string) error {
	rule, err := r.Rule(path, name, methods)
	if err != nil {
		return err
	}

	parts := strings.Split(name, ".")
	method, ok := context.MethodByName(parts[len(parts)-1])
	if !ok {
		return nil // change
	}

	r.endpoints[rule] = &endpoint{
		context:  context,
		method:   method.Func,
		bindings: make(map[string]interface{}),
	}

	return nil
}

// Dispatch returns DispatchFunc that waitress/router expects. The DispatchFunc
// performs route matching and constructs the context for the endpoint's method
// receiver.
func (r *Router) Dispatch(ctx *Context) router.DispatchFunc {
	return func(rule *router.Rule, args map[string]interface{}) interface{} {
		// Find the endpoint given the matched rule.
		endpoint, ok := r.endpoints[rule]
		if !ok {
			return r.NotFoundHandler()
		}

		// Ensure that arguments provided match the number of arguments expected.
		t := endpoint.method.Type()
		keys := rule.Parameters()
		if t.NumIn() > len(keys)+1 {
			return r.InternalServerErrorHandler()
		}

		// Construct the method receiver for the endpoint.
		receiver := reflect.New(endpoint.context.Elem())
		receiver.Elem().FieldByName("Context").Set(reflect.ValueOf(ctx))
		for name, value := range endpoint.bindings {
			receiver.Elem().FieldByName(name).Set(reflect.ValueOf(value))
		}

		// Prepare the calling parameters.
		// Method expressions take the receiver as the first argument.
		params := make([]reflect.Value, len(keys)+1)
		params[0] = receiver
		for i, key := range keys {
			params[i+1] = reflect.ValueOf(args[key])
		}

		// Call our endpoint and return successful if no return value.
		rv := endpoint.method.Call(params)
		if len(rv) == 0 {
			return []byte(nil)
		}

		// We do not support more than one return value.
		if len(rv) > 1 {
			return r.InternalServerErrorHandler()
		}

		return rv[0].Interface()
	}
}

// ServeHTTP implements the http.Handler interface. It binds the router to the
// current request, creates the context, and dispatches to the DispatchFunc.
// The response can be a byte slice, a string, an http.Handler, or any function
// with the method signature of an http.HandlerFunc. If the return value is
// anything else, the InternalServerErrorHandler will be invoked.
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	adapter := r.BindToRequest(req)

	w := NewResponseWriter(rw)
	ctx := NewContext(w, req, adapter)

	rv := adapter.Dispatch(r.Dispatch(ctx))

	switch v := rv.(type) {
	case []byte:
		w.Write(v)
	case string:
		w.Write([]byte(v))
	case http.Handler:
		v.ServeHTTP(w, req)
	case func(http.ResponseWriter, *http.Request):
		v(w, req)
	default:
		fallback := r.InternalServerErrorHandler()
		fallback.ServeHTTP(w, req)
	}
}
