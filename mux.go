package gohttputil

import (
	"net/http"
	"slices"

	"github.com/rs/cors"
)

// Mux is a wrapper around the http.ServeMux.
// It provides a http.Handler implementation that
// basically calls http.ServeMux.ServeHTTP with or
// without CORS wrapper.
//
// Any route handler attached through this Mux adds that route
// to the internal http.ServeMux by wrapping middlewares
// from different levels (global, group level or route level).
//
// Global middlewares are added by Mux.Use method. These middlewares
// are applied to all routes defined by this instance of the Mux.
//
// Group level middlewares are only applied to the group routes.
//
// Route level middlewares are applied per route per method.
type Mux struct {
	mux         *http.ServeMux
	middlewares []Middleware
	corsHandler *cors.Cors
}

// New creates a new instance of Mux.
func New() *Mux {
	return &Mux{
		mux:         &http.ServeMux{},
		middlewares: []Middleware{},
	}
}

// ServeHTTP implements http.Handler.
// This calls the internal http.ServeMux.ServeHTTP with
// or without CORS wrapper handler.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.corsHandler != nil {
		m.corsHandler.ServeHTTP(w, r, m.mux.ServeHTTP)
		return
	}
	m.mux.ServeHTTP(w, r)
}

// Use appends middlewares to global middlewares.
func (r *Mux) Use(middlewares ...Middleware) *Mux {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

// Route implements Router.
func (r *Mux) Route(route string) RouteHandler {
	return &routeHandler{
		mux:             r.mux,
		route:           route,
		rootMiddlewares: slices.Clone(r.middlewares),
		middlewares:     []Middleware{},
	}
}

// Group implements Grouper.
func (m *Mux) Group(prefix string) Group {
	return &group{
		mux:         m.mux,
		prefix:      prefix,
		middlewares: slices.Clone(m.middlewares),
	}
}

// EnableCORS wraps the internal http.ServeMux with CORS handler.
// Without any option in argument, it allows all methods, origins and
// headers.
func (m *Mux) EnableCORS(opt ...cors.Options) {
	c := cors.AllowAll()
	if len(opt) > 0 {
		c = cors.New(opt[0])
	}
	m.corsHandler = c
}

var (
	_ = (http.Handler)(&Mux{})
	_ = (Router)(&Mux{})
	_ = (Grouper)(&Mux{})
)
