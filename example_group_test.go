package gohttputil_test

import gohttputil "github.com/asif-mahmud/go-httputil"

// Example_group shows grouping multiple routes
func Example_group() {
	// create mux instance
	m := gohttputil.New()

	// define a group of routes
	m.Group("/api/v1").
		// group level middlewares. these middlewares
		// will be applied to all routes defined in this group
		Use(groupMiddleware1).
		Use(groupMiddleware2).

		// define handlers for a route
		Route("/orders", func(rh gohttputil.RouteHandler) {
			rh.
				// no route level middleware will be applied to GET handler
				Get(handler1).

				// middleware1 will be applied to POST handler
				Use(middleware1).
				Post(handler1).

				// no middleware applied to PUT method handler
				Put(handler1)
		}).

		// define handlers for another route under same group
		Route("/users", func(rh gohttputil.RouteHandler) {
			rh.
				// no route level middleware will be applied to GET handler
				Get(handler1).

				// middleware1 will be applied to POST handler
				Use(middleware1).
				Post(handler1).

				// no middleware applied to PUT method handler
				Put(handler1)
		})
}
