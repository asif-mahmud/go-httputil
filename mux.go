package gohttputil

import (
	"net/http"
	"slices"
)

// Mux is a wrapper around the http.ServeMux.
// It provides a http.Handler implementation that
// basically calls http.DefaultServeMux.ServeHTTP.
//
// Any route handler attached through this Mux adds that route
// to the http.DefaultServeMux by wrapping middlewares
// from different levels (global, group level or route level).
//
// Global middlewares are added by Mux.Use method. These middlewares
// are applied to all routes defined by this instance of the Mux.
//
// Group level middlewares are only applied to the group routes.
//
// Route level middlewares are applied per route per method.
type Mux struct {
	middlewares []Middleware
}

// New creates a new instance of Mux.
// Multiple instances of the Mux doesn't isolate their routes
// as the routes are directly attached to the http.DefaultServeMux.
func New() *Mux {
	return &Mux{
		middlewares: []Middleware{},
	}
}

// ServeHTTP implements http.Handler.
// This just calls http.DefaultServeMux.ServeHTTP.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}

// Use appends middlewares to global middlewares.
func (r *Mux) Use(middlewares ...Middleware) *Mux {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

// Route implements Router.
func (r *Mux) Route(route string) RouteHandler {
	return &routeHandler{
		route:           route,
		rootMiddlewares: slices.Clone(r.middlewares),
		middlewares:     []Middleware{},
	}
}

// Group implements Grouper.
func (m *Mux) Group(prefix string) Group {
	return &group{
		prefix:      prefix,
		middlewares: slices.Clone(m.middlewares),
	}
}

var (
	_ = (http.Handler)(&Mux{})
	_ = (Router)(&Mux{})
	_ = (Grouper)(&Mux{})
)
