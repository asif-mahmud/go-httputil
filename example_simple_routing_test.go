package gohttputil_test

import (
	"log"
	"net/http"

	gohttputil "github.com/asif-mahmud/go-httputil"
)

// Example_simpleRouting shows defining routes
func Example_simpleRouting() {
	// create mux instance
	m := gohttputil.New()

	// add global middlewares.
	// these middlewares will be applied to all routes.
	m.
		Use(globalMiddleware1).
		Use(globalMiddleware2)

	// define routes
	m.Route("/api/order").
		// no middleware applied to GET method handler
		Get(handler1).

		// middleware1 and middleware2 are applied to POST method handler
		Use(middleware1).
		Use(middleware2).
		Post(handler2).

		// no middleware applied to PUT method handler
		Put(handler1).

		// only middleware1 applied to DELETE method handler
		Use(middleware1).
		Delete(handler2)

	// you can use the created mux now
	log.Fatal(http.ListenAndServe(":3000", m))

	// or you can use http.Server with handler set to the mux
	server := http.Server{
		Handler: m,
		Addr:    ":3000",
	}
	log.Fatal(server.ListenAndServe())
}
