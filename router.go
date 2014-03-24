package waitress

import (
	"net/http"
	"reflect"

	"github.com/pnelson/waitress/router"
)

type Router struct {
	*router.Router
	context   reflect.Type
	endpoints map[*router.Rule]reflect.Value
}

func NewRouter(ctx interface{}) *Router {
	return &Router{
		Router:    router.New(),
		context:   reflect.TypeOf(ctx),
		endpoints: make(map[*router.Rule]reflect.Value),
	}
}

func (r *Router) Route(path, name string, methods []string) error {
	rule, err := r.Rule(path, name, methods)
	if err != nil {
		return err
	}

	endpoint, ok := r.context.MethodByName(name)
	if !ok {
		return nil // change
	}

	r.endpoints[rule] = endpoint.Func

	return nil
}

func (r *Router) Mount(prefix string, fragment *Fragment) {
	//name := reflect.TypeOf(fragment).Elem().Name()
}

func (r *Router) Dispatch(w http.ResponseWriter, req *http.Request) router.DispatchFunc {
	// Construct the method receiver for the endpoint.
	ctx := reflect.New(r.context.Elem())
	ctx.Elem().FieldByName("Context").Set(reflect.ValueOf(NewContext(w, req)))

	return func(rule *router.Rule, args map[string]interface{}) interface{} {
		// Find the endpoint given the matched rule.
		endpoint, ok := r.endpoints[rule]
		if !ok {
			return r.NotFoundHandler()
		}

		// Ensure that arguments provided match the number of arguments expected.
		t := endpoint.Type()
		keys := rule.Parameters()
		if t.NumIn() > len(keys)+1 {
			return r.InternalServerErrorHandler()
		}

		// Prepare the calling parameters.
		// Method expressions take the receiver as the first argument.
		params := make([]reflect.Value, len(keys)+1)
		params[0] = ctx
		for i, key := range keys {
			params[i+1] = reflect.ValueOf(args[key])
		}

		// Call our endpoint and return successful if no return value.
		rv := endpoint.Call(params)
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

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	adapter := r.BindToRequest(req)
	rv := adapter.Dispatch(r.Dispatch(w, req))

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
