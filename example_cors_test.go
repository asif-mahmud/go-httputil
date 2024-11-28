package gohttputil_test

import (
	"log"
	"net/http"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/rs/cors"
)

// Example_cors sets up CORS preflight handler globally.
func Example_cors() {
	// create mux instance
	m := gohttputil.New()

	// enable CORS handler globally with default
	// all allowed option
	m.EnableCORS()

	// or enable CORS handler with your own option
	m.EnableCORS(cors.Options{})

	// you can use the created mux now
	log.Fatal(http.ListenAndServe(":3000", m))

	// or you can use http.Server with handler set to the mux
	server := http.Server{
		Handler: m,
		Addr:    ":3000",
	}
	log.Fatal(server.ListenAndServe())
}
