package otters

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	Mux  *http.ServeMux
	port string
	host string

	middleware         []Middleware
	settings           Settings
	renderer           Renderer
	settingsConfigured bool
}

type Middleware func(http.Handler) http.Handler

// New creates a new default otters application
func New() App {
	mux := http.NewServeMux()
	settings := DefaultSettings()
	return App{
		Mux:                mux,
		port:               "8008",
		host:               "",
		middleware:         []Middleware{},
		settings:           settings,
		renderer:           NewHTMLRenderer(settings),
		settingsConfigured: false,
	}
}

// SubApp creates a child app that inherits settings + renderer from parent
// Caller should mount it via parent.MountApp(...)
func SubApp(parent *App) App {
	mux := http.NewServeMux()
	return App{
		Mux:        mux,
		port:       parent.port,
		host:       parent.host,
		middleware: []Middleware{},
		settings:   parent.settings,
		renderer:   parent.renderer,
	}
}

// Handler returns the app's mux wrapped in app-level middleware.
func (a *App) Handler() http.Handler {
	var h http.Handler = a.Mux
	for i := len(a.middleware) - 1; i >= 0; i-- {
		h = a.middleware[i](h)
	}
	return h
}

// Mount mountes an http.Handler under prefix.
func (a *App) Mount(prefix string, h http.Handler) {
	prefix = normalizeMountPrefix(prefix)

	pattern := prefix + "/"
	a.Mux.Handle(pattern, http.StripPrefix(prefix, h))
}

// MountApp mounts another otters app under the prefix.
func (a *App) MountApp(prefix string, child *App) {
	a.Mount(prefix, child.Handler())
}

// GetNetworkLocation returns the network location that the app will listen on.
func (a *App) GetNetworkLocation() string {
	return fmt.Sprintf("%s:%s", a.host, a.port)
}

// Serve sets the app to listen on the set network location and return any error.
func (a *App) Serve() error {
	handler := http.Handler(a.Mux)
	for i := len(a.middleware) - 1; i >= 0; i-- {
		handler = a.middleware[i](handler)
	}

	fmt.Printf("listening on %s\n", a.GetNetworkLocation())
	return http.ListenAndServe(a.GetNetworkLocation(), handler)
}

// Use registers a middleware handler to all routes of the application
func (a *App) Use(m Middleware) {
	a.middleware = append(a.middleware, m)
}

// HandleFunc calls out to the underlying net/http HandleFunc
func (a *App) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	a.Mux.HandleFunc(pattern, handler)
}

// Handle calls out to the underlying net/http Handle
func (a *App) Handle(pattern string, handler http.Handler) {
	a.Mux.Handle(pattern, handler)
}

// Get creates a get route for the application
// pattern: the route to listen to.
// handler: takes the otters Context and writes to it.
// Ex:
//
//	app.Get("/ping", func(ctx otters.Ctx) {
//	   ctx.Write("ok")
//	})
func (a *App) Get(pattern string, handler func(*Ctx), middleware ...Middleware) {
	route := fmt.Sprintf("GET %s", pattern)
	funcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := a.newCtx(w, r)
		handler(ctx)
	})

	finalHandler := applyMiddleware(funcHandler, middleware...)
	a.Mux.Handle(route, finalHandler)
}

// Post creates a POST request route for the application
//
// pattern: the route to listen to
//
// handler: takes the otters Context and writes to it
func (a *App) Post(pattern string, handler func(*Ctx), middleware ...Middleware) {
	route := fmt.Sprintf("POST %s", pattern)
	funcHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := a.newCtx(w, r)
		handler(ctx)
	})
	finalHandler := applyMiddleware(funcHandler, middleware...)
	a.Mux.Handle(route, finalHandler)
}

// WithStatic tells the app to serve files from the directory on the given route
func (a *App) WithStatic(pattern, dir string) {
	fs := http.FileServer(http.Dir(dir))
	a.Mux.Handle(pattern, http.StripPrefix(pattern, fs))
}

// WithPort sets the port for the application to listen on
func (a *App) WithPort(port string) {
	a.port = port
}

// WithHost sets the host for the application to listen on
func (a *App) WithHost(host string) {
	a.host = host
}

// WithSettings sets the application settings and reconfigures dependent settings
func (a *App) WithSettings(settings Settings) {
	a.settings = settings
	a.renderer = NewHTMLRenderer(settings)

	if !a.settingsConfigured && settings.Static.Route != "" && settings.Static.Dir != "" {
		a.WithStatic(settings.Static.Route, settings.Static.Dir)
		a.settingsConfigured = true
	}
}

func applyMiddleware(handler http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

func (a *App) newCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		Writer:   w,
		Request:  r,
		settings: a.settings,
		html:     a.renderer,
	}
}

func normalizeMountPrefix(prefix string) string {
	if prefix == "" {
		return "/"
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if prefix != "/" && strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}
	return prefix
}
