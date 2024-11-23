package gohttputil

import "slices"

type Grouper interface {
	Group(string) Group
}

type GroupRouter func(RouteHandler)

type Group interface {
	Use(...Middleware) Group
	Route(string, GroupRouter) Group
}

type group struct {
	prefix      string
	middlewares []Middleware
}

// Route implements Group.
func (g *group) Route(route string, router GroupRouter) Group {
	router(&routeHandler{
		route:           g.prefix + route,
		rootMiddlewares: slices.Clone(g.middlewares),
		middlewares:     []Middleware{},
	})

	return g
}

// Use implements Group.
func (g *group) Use(middlewares ...Middleware) Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

var _ = (Group)(&group{})
