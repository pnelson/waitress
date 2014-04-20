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

	app.RedirectHandler = func(path string, code int) http.Handler {
		return RedirectToWithCode(path, code)
	}

	app.NotFoundHandler = func() http.Handler {
		return NotFound()
	}

	app.MethodNotAllowedHandler = func(allowed []string) http.Handler {
		return MethodNotAllowed(allowed)
	}

	app.InternalServerErrorHandler = func() http.Handler {
		return InternalServerError()
	}

	return app
}

func (app *Application) Dispatch(w http.ResponseWriter, r *http.Request) {
	defer app.Recover(w, r)
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

func (app *Application) Recover(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		handler := app.InternalServerErrorHandler()
		handler.ServeHTTP(w, r)
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
