package gohttputil

import "net/http"

type Mux struct {
	globalMiddlewares []Middleware
}

var (
	_ = (http.Handler)(&Mux{})
	_ = (Router)(&Mux{})
)

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
		route:       route,
		mux:         r,
		middlewares: []Middleware{},
	}
}
