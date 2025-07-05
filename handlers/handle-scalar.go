package handlers

import (
	"embed"
	"io"
	"net/http"
)

//go:embed scalar-dist/*
var scalarDistFS embed.FS

// HandleScalar returns a handler function to serve api documentation via scalar ui.
//
// Scalar website - https://scalar.com/
// Scalar documentation - https://guides.scalar.com/scalar/introduction
//
// This has a more modern Web UI for OpenAPI browser based documentation.
// To attach this handler to a path do this -
//
// mux.Route("/swagger/{path...}").Get(HandleScalar(doc, "path"))
//
// This will let the handler server swagger.json and other static files
// under /swagger path.
func HandleScalar(doc io.Reader, pathKey string) http.HandlerFunc {
	return swaggerHandlerFactory(doc, pathKey, scalarDistFS, "scalar-dist")
}
