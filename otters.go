package otters

import (
	"fmt"
	"net/http"
	"time"
)

type App struct {
	Mux  *http.ServeMux
	Port string
	Host string

	middleware []Middleware
}

type Middleware func(http.Handler) http.Handler

type Ctx struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

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
func (a *App) Get(pattern string, handler func(Ctx)) {
	route := fmt.Sprintf("GET %s", pattern)
	a.Mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ctx := Ctx{
			Writer:  w,
			Request: r,
		}
		handler(ctx)
	})
}

// Create a post request route for the application
//
// pattern: the route to listen to
//
// handler: takes the otters Context, and writes to it
func (a *App) Post(pattern string, handler func(Ctx)) {
	route := fmt.Sprintf("POST %s", pattern)
	a.Mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ctx := Ctx{
			Writer:  w,
			Request: r,
		}
		handler(ctx)
	})
}

// Tells the app to serve files from the directory on the given route
func (a *App) WithStatic(pattern, dir string) {
	fs := http.FileServer(http.Dir(dir))
	a.Mux.Handle(pattern, http.StripPrefix(pattern, fs))
}

func (a *App) WithDefaultLogger() {
	logger := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(ww, r)
			duration := time.Since(start)
			fmt.Printf("%s %s %v %d\n", r.Method, r.URL.Path, duration, ww.statusCode)
		})
	}

	a.Middleware(logger)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Writes the given string to the otters Context
func (c Ctx) String(content string) {
	fmt.Fprint(c.Writer, content)
}
