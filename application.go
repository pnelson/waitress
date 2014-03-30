package waitress

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/pnelson/waitress/middleware"
)

type Application struct {
	*middleware.Builder
	*Router

	context reflect.Type
	closed  bool
}

func New(ctx interface{}) *Application {
	app := &Application{
		Builder: &middleware.Builder{},
		Router:  NewRouter(),
		context: reflect.TypeOf(ctx),
	}

	app.SetRedirectHandler(func(path string, code int) http.Handler {
		return RedirectToWithCode(path, code)
	})

	app.SetNotFoundHandler(func() http.Handler {
		return NotFound()
	})

	app.SetMethodNotAllowedHandler(func(allowed []string) http.Handler {
		return MethodNotAllowed(allowed)
	})

	app.SetInternalServerErrorHandler(func() http.Handler {
		return InternalServerError()
	})

	return app
}

func (app *Application) Dispatch(w http.ResponseWriter, r *http.Request) {
	defer app.Recover()
	if !app.closed {
		app.UseHandler(app.Router)
		app.closed = true
	}
	app.Builder.ServeHTTP(w, r)
}

func (app *Application) Route(path, name string, methods []string) error {
	return app.Router.Route(path, name, app.context, methods)
}

func (app *Application) Mount(prefix, name string, fragment *Fragment) error {
	return fragment.Register(app, prefix, name)
}

func (app *Application) Recover() {
	if err := recover(); err != nil {
		fmt.Println("recovered from panic:", err)
	}
}

func (app *Application) Run() {
	addr := fmt.Sprintf("%s:%d", "localhost", 3000)
	fmt.Println(fmt.Sprintf("Running on %s://%s/", "http", addr))

	http.ListenAndServe(addr, app)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Dispatch(w, r)
}

func (app *Application) SetRedirectHandler(f func(string, int) http.Handler) {
	app.Router.RedirectHandler = f
}

func (app *Application) SetNotFoundHandler(f func() http.Handler) {
	app.Router.NotFoundHandler = f
}

func (app *Application) SetMethodNotAllowedHandler(f func([]string) http.Handler) {
	app.Router.MethodNotAllowedHandler = f
}

func (app *Application) SetInternalServerErrorHandler(f func() http.Handler) {
	app.Router.InternalServerErrorHandler = f
}
