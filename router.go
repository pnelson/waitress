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
	ctx := reflect.New(r.context.Elem())
	ctx.Elem().FieldByName("Context").Set(reflect.ValueOf(NewContext(w, req)))

	return func(rule *router.Rule, args map[string]interface{}) (interface{}, error) {
		endpoint, ok := r.endpoints[rule]
		if !ok {
			return nil, nil // change
		}

		t := endpoint.Type()
		keys := rule.Parameters()
		if t.NumIn() > len(keys)+1 {
			return nil, nil // change
		}

		params := make([]reflect.Value, len(keys)+1)
		params[0] = ctx
		for i, key := range keys {
			params[i+1] = reflect.ValueOf(args[key])
		}

		rv := endpoint.Call(params)
		if count := len(rv); count != 1 {
			if count == 0 {
				return nil, nil
			}
			return nil, nil // change
		}

		return rv[0].Interface(), nil
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	adapter := r.BindToRequest(req)
	rv, err := adapter.Dispatch(r.Dispatch(w, req))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	switch response := rv.(type) {
	case []byte:
		w.Write(response)
	case string:
		w.Write([]byte(response))
	case http.Handler:
		response.ServeHTTP(w, req)
	case http.HandlerFunc:
		response(w, req)
	default:
		http.Error(w, "Internal Server Error", 500)
	}
}
