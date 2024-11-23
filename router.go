package gohttputil

import (
	"fmt"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)

type Router interface {
	Route(string) RouteHandler
}

type RouteHandler interface {
	Use(...Middleware) RouteHandler

	Get(Handler) RouteHandler
	Post(Handler) RouteHandler
	Put(Handler) RouteHandler
	Patch(Handler) RouteHandler
	Delete(Handler) RouteHandler
}

type routeHandler struct {
	route           string
	rootMiddlewares []Middleware
	middlewares     []Middleware
}

// Use implements Router.
func (r *routeHandler) Use(middlewares ...Middleware) RouteHandler {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *routeHandler) createHandler(method string, handler Handler) {
	var h http.Handler
	h = http.HandlerFunc(handler)

	for _, m := range r.rootMiddlewares {
		h = m(h)
	}

	for _, m := range r.middlewares {
		h = m(h)
	}

	http.Handle(fmt.Sprintf("%s %s", method, r.route), h)
}

func newRouter(r *routeHandler) *routeHandler {
	return &routeHandler{
		route:           r.route,
		rootMiddlewares: r.rootMiddlewares,
		middlewares:     []Middleware{},
	}
}

// Get implements RouteHandler.
func (r *routeHandler) Get(handler Handler) RouteHandler {
	r.createHandler(http.MethodGet, handler)
	return newRouter(r)
}

// Post implements RouteHandler.
func (r *routeHandler) Post(handler Handler) RouteHandler {
	r.createHandler(http.MethodPost, handler)
	return newRouter(r)
}

// Put implements RouteHandler.
func (r *routeHandler) Put(handler Handler) RouteHandler {
	r.createHandler(http.MethodPut, handler)
	return newRouter(r)
}

// Patch implements RouteHandler.
func (r *routeHandler) Patch(handler Handler) RouteHandler {
	r.createHandler(http.MethodPatch, handler)
	return newRouter(r)
}

// Delete implements RouteHandler.
func (r *routeHandler) Delete(handler Handler) RouteHandler {
	r.createHandler(http.MethodDelete, handler)
	return newRouter(r)
}

var _ = (RouteHandler)(&routeHandler{})
