package handlers

import (
	"embed"
	"io"
	"net/http"
)

//go:embed swagger-dist/*
var distFS embed.FS

// HandleSwagger returns a handler function to serve swagger doc and related files.
//
// To attach this handler to a path do this -
//
// mux.Route("/swagger/{path...}").Get(HandleSwagger(doc, "path"))
//
// This will let the handler serve swagger.json and other static files
// under /swagger path.
func HandleSwagger(doc io.Reader, pathKey string) http.HandlerFunc {
	return swaggerHandlerFactory(doc, pathKey, distFS, "swagger-dist")
}
