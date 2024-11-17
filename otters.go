package otters

import (
	"fmt"
	"net/http"
)

type App struct {
	Mux  *http.ServeMux
	Port string
	Host string
}

type Ctx struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func New() App {
	mux := http.NewServeMux()
	return App{
		Mux:  mux,
		Port: "8080",
		Host: "",
	}
}

// Returns the network location that the app will listen on.
func (a *App) GetNetworkLocation() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

// Sets the app to listen on the set network location and return any error.
func (a *App) Serve() error {
	fmt.Printf("listening on %s\n", a.GetNetworkLocation())
	return http.ListenAndServe(a.GetNetworkLocation(), a.Mux)
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

// Writes the given string to the otters Context
func (c Ctx) String(content string) {
	fmt.Fprint(c.Writer, content)
}
