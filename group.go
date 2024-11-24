package gohttputil

import "slices"

// Grouper defines interface to create a new Group.
type Grouper interface {
	Group(string) Group
}

// GroupRouter defines function signature for group level router.
type GroupRouter func(RouteHandler)

// Group defines group level routing methods.
type Group interface {
	// Use adds middlewares to group level middleware list.
	// These middlewares are applied to all routes defined under this group.
	// This acts like global middlewares but applied only for this group.
	Use(...Middleware) Group

	// Route creates a new RouteHandler under the current Group and calls GroupRouter with it.
	// This lets user to define one or more routes under the current Group.
	Route(string, GroupRouter) Group
}

type group struct {
	prefix      string
	middlewares []Middleware
}

// Use implements Group.
func (g *group) Use(middlewares ...Middleware) Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
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

var _ = (Group)(&group{})
