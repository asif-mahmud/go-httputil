package gohttputil

import (
	"fmt"
	"net/http"
	"slices"
)

// Handler
type Handler func(http.ResponseWriter, *http.Request)

// Router interface to create a RouteHandler
type Router interface {
	Route(string) RouteHandler
}

// RouteHandler interface to configure route handlers.
type RouteHandler interface {
	// Use adds middlewares to be used for the current route and current method.
	// This should be called before calling any of the http method handlers (Get,
	// Post, Put, Patch or Delete method). After calling an http method handler
	// this list of middlewares are cleared so that user can define different set
	// of middlewares for next http method handler. So route level middlewares
	// are defined per route per http method.
	Use(...Middleware) RouteHandler

	// Get attaches handler to http GET method
	Get(Handler) RouteHandler

	// Post attaches handler to http POST method
	Post(Handler) RouteHandler

	// Put attaches handler to http PUT method
	Put(Handler) RouteHandler

	// Patch attaches handler to http PATCH method
	Patch(Handler) RouteHandler

	// Delete attaches handler to http DELETE method
	Delete(Handler) RouteHandler
}

type routeHandler struct {
	route           string
	rootMiddlewares []Middleware
	middlewares     []Middleware
}

// Use implements RouteHandler.
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
		rootMiddlewares: slices.Clone(r.rootMiddlewares),
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
