package gohttputil

import (
	"fmt"
	"net/http"
	"slices"
)

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
	// of middlewares for next http method handler. Soel middlewares
	// are defined per route per http method.
	Use(...Middleware) RouteHandler

	// Get attaches handler to http GET method
	Get(http.HandlerFunc) RouteHandler

	// Post attaches handler to http POST method
	Post(http.HandlerFunc) RouteHandler

	// Put attaches handler to http PUT method
	Put(http.HandlerFunc) RouteHandler

	// Patch attaches handler to http PATCH method
	Patch(http.HandlerFunc) RouteHandler

	// Delete attaches handler to http DELETE method
	Delete(http.HandlerFunc) RouteHandler
}

type routeHandler struct {
	mux             *http.ServeMux
	route           string
	rootMiddlewares []Middleware
	middlewares     []Middleware
}

// Use implements RouteHandler.
func (r *routeHandler) Use(middlewares ...Middleware) RouteHandler {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *routeHandler) createHandler(method string, handler http.HandlerFunc) {
	var h http.Handler
	h = http.HandlerFunc(handler)

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	for i := len(r.rootMiddlewares) - 1; i >= 0; i-- {
		h = r.rootMiddlewares[i](h)
	}

	r.mux.Handle(fmt.Sprintf("%s %s", method, r.route), h)
}

func (r *routeHandler) reset() RouteHandler {
	r.middlewares = []Middleware{}
	r.rootMiddlewares = slices.Clone(r.rootMiddlewares)
	return r
}

// Get implements RouteHandler.
func (r *routeHandler) Get(handler http.HandlerFunc) RouteHandler {
	r.createHandler(http.MethodGet, handler)
	return r.reset()
}

// Post implements RouteHandler.
func (r *routeHandler) Post(handler http.HandlerFunc) RouteHandler {
	r.createHandler(http.MethodPost, handler)
	return r.reset()
}

// Put implements RouteHandler.
func (r *routeHandler) Put(handler http.HandlerFunc) RouteHandler {
	r.createHandler(http.MethodPut, handler)
	return r.reset()
}

// Patch implements RouteHandler.
func (r *routeHandler) Patch(handler http.HandlerFunc) RouteHandler {
	r.createHandler(http.MethodPatch, handler)
	return r.reset()
}

// Delete implements RouteHandler.
func (r *routeHandler) Delete(handler http.HandlerFunc) RouteHandler {
	r.createHandler(http.MethodDelete, handler)
	return r.reset()
}

var _ = (RouteHandler)(&routeHandler{})
