/*
Package waitress is an HTTP framework for Go.

The framework and the sub-packages attempt to stay close to net/http but
provide abstractions over actions that are performed frequently.
*/
package waitress

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"

	"github.com/pnelson/waitress/middleware"
)

// Application holds the top level configuration.
type Application struct {
	*middleware.Builder // The middleware builder is embedded for its methods.
	*Router             // The router is embedded for its methods.

	context reflect.Type
	closed  bool
}

// New returns a new Application preconfigured with some sane defaults.
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

// Dispatch kicks off the middleware processing for each request.
func (app *Application) Dispatch(w http.ResponseWriter, r *http.Request) {
	defer app.Recover(w, r)
	if !app.closed {
		app.UseHandler(app.Router)
		app.closed = true
	}
	app.Builder.ServeHTTP(w, r)
}

// Route registers a route with the application context.
func (app *Application) Route(path, name string, methods []string) error {
	return app.Router.Route(path, name, app.context, methods)
}

// Mount registers a Fragment with the Application at a defined prefix.
func (app *Application) Mount(prefix, name string, fragment *Fragment) error {
	return fragment.Register(app, prefix, name)
}

// Recover will dump a stack trace and process the InternalServerErrorHandler
// on panic.
func (app *Application) Recover(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		log.Println(err)
		debug.PrintStack()
		handler := app.InternalServerErrorHandler()
		handler.ServeHTTP(w, r)
	}
}

// Run will serve the application on localhost:3000 for now. Of course, this
// Application object is just another http.Handler, so you can ListenAndServe
// on your own.
func (app *Application) Run() {
	addr := fmt.Sprintf("%s:%d", "localhost", 3000)
	fmt.Println(fmt.Sprintf("Running on %s://%s/", "http", addr))

	http.ListenAndServe(addr, app)
}

// ServeHTTP implements the http.Handler interface.
func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Dispatch(w, r)
}
