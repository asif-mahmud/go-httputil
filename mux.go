package gohttputil

import (
	"net/http"
	"slices"
)

type Mux struct {
	globalMiddlewares []Middleware
}

// Group implements Grouper.
func (m *Mux) Group(prefix string) Group {
	return &group{
		prefix:      prefix,
		middlewares: slices.Clone(m.globalMiddlewares),
	}
}

func New() *Mux {
	return &Mux{
		globalMiddlewares: []Middleware{},
	}
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}

func (r *Mux) Use(middlewares ...Middleware) *Mux {
	r.globalMiddlewares = append(r.globalMiddlewares, middlewares...)
	return r
}

func (r *Mux) Route(route string) RouteHandler {
	return &routeHandler{
		route:           route,
		rootMiddlewares: r.globalMiddlewares,
		middlewares:     []Middleware{},
	}
}

var (
	_ = (http.Handler)(&Mux{})
	_ = (Router)(&Mux{})
	_ = (Grouper)(&Mux{})
)
