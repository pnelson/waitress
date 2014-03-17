package waitress

import (
	"fmt"
	"net/http"

	"github.com/pnelson/waitress/middleware"
)

type Application struct {
	*middleware.Builder
	*Router
}

func New() *Application {
	return &Application{
		Builder: &middleware.Builder{},
		Router:  NewRouter(),
	}
}

func (app *Application) Close() {
	app.UseHandler(app.Router)
}

func (app *Application) Dispatch(w http.ResponseWriter, r *http.Request) {
	defer app.Recover()
	app.Builder.ServeHTTP(w, r)
}

func (app *Application) Recover() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}

func (app *Application) Run() {
	app.Close()

	addr := fmt.Sprintf("%s:%d", "localhost", 3000)
	fmt.Println(fmt.Sprintf("Running on %s://%s/", "http", addr))

	http.ListenAndServe(addr, app)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Dispatch(w, r)
}
