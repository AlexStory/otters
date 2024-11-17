package otters

import (
	"fmt"
	"net/http"
)

type App struct {
	Mux  *http.ServeMux
	Port string
	Host string

	middleware []Middleware
}

type Middleware func(http.Handler) http.Handler

// Creates a new default otters application
func New() App {
	mux := http.NewServeMux()
	return App{
		Mux:        mux,
		Port:       "8008",
		Host:       "",
		middleware: []Middleware{},
	}
}

// Returns the network location that the app will listen on.
func (a *App) GetNetworkLocation() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

// Sets the app to listen on the set network location and return any error.
func (a *App) Serve() error {
	handler := http.Handler(a.Mux)
	for i := len(a.middleware) - 1; i >= 0; i-- {
		handler = a.middleware[i](handler)
	}

	fmt.Printf("listening on %s\n", a.GetNetworkLocation())
	return http.ListenAndServe(a.GetNetworkLocation(), handler)
}

// Adds a middleware handler to all routes of the application
func (a *App) Middleware(m Middleware) {
	a.middleware = append(a.middleware, m)
}

// Calls out to the underlying net/http HandleFunc
func (a *App) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	a.Mux.HandleFunc(pattern, handler)
}

// Calls out to the underlying net/http Handle
func (a *App) Handle(pattern string, handler http.Handler) {
	a.Mux.Handle(pattern, handler)
}

// Create a get route for the appplication
// pattern: the route to listen to.
// handler takes the otters Context, and writes to it.
// Ex:
//
//	app.Get("/ping", func(ctx otters.Ctx) {
//	   ctx.Write("ok")
//	})
func (a *App) Get(pattern string, handler func(Ctx), middleware ...Middleware) {
	route := fmt.Sprintf("GET %s", pattern)
	funcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := Ctx{
			Writer:  w,
			Request: r,
		}
		handler(ctx)
	})

	finalHandler := applyMiddleware(funcHandler, middleware...)
	a.Mux.Handle(route, finalHandler)
}

// Create a post request route for the application
//
// pattern: the route to listen to
//
// handler: takes the otters Context, and writes to it
func (a *App) Post(pattern string, handler func(Ctx), middleware ...Middleware) {
	route := fmt.Sprintf("POST %s", pattern)
	funcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := Ctx{
			Writer:  w,
			Request: r,
		}
		handler(ctx)
	})
	finalHandler := applyMiddleware(funcHandler, middleware...)
	a.Mux.Handle(route, finalHandler)
}

// Tells the app to serve files from the directory on the given route
func (a *App) WithStatic(pattern, dir string) {
	fs := http.FileServer(http.Dir(dir))
	a.Mux.Handle(pattern, http.StripPrefix(pattern, fs))
}

func applyMiddleware(handler http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}
