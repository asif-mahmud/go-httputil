package handlers

import (
	"bytes"
	"embed"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/asif-mahmud/go-httputil/helpers"
	golog "github.com/asif-mahmud/go-log"
)

//go:embed dist/*
var distFS embed.FS

var swaggerDoc []byte

// HandleSwagger returns a handler function to serve swagger doc and related files.
//
// To attach this handler to a path do this -
//
// mux.Route("/swagger/{path...}").Get(HandleSwagger(doc, "path"))
func HandleSwagger(doc io.Reader, pathKey string) http.HandlerFunc {
	data, err := io.ReadAll(doc)
	if err != nil {
		slog.Error("Failed to load swagger.json file", golog.Extra(map[string]any{
			"error": err.Error(),
		}))
	} else {
		swaggerDoc = bytes.Clone(data)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		filePath := r.PathValue(pathKey)

		// if file path not specified serve index file
		if filePath == "" || filePath == "/" {
			filePath = "index.html"
		}

		// for swagger json file
		if strings.HasSuffix(filePath, "swagger.json") {
			if swaggerDoc == nil || len(swaggerDoc) == 0 {
				helpers.SendError(w, http.StatusNotFound, "File not found", nil)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(swaggerDoc)
			return
		}

		// for static dist files
		data, err := distFS.ReadFile(path.Join("dist", filePath))
		if err != nil {
			helpers.SendError(w, http.StatusNotFound, "File not found", nil)
			return
		}
		ext := path.Ext(filePath)
		mime := mime.TypeByExtension(ext)
		w.Header().Add("Content-Type", mime)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

	return fn
}
