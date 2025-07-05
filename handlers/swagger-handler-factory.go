package handlers

import (
	"bytes"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/asif-mahmud/go-httputil/helpers"
	golog "github.com/asif-mahmud/go-log"
)

func swaggerHandlerFactory(
	doc io.Reader,
	pathKey string,
	fsys fs.FS,
	fsRootDir string,
) http.HandlerFunc {
	data, err := io.ReadAll(doc)
	docData := []byte{}
	if err != nil {
		slog.Error("Failed to load swagger.json file", golog.Extra(map[string]any{
			"error": err.Error(),
		}))
	} else {
		docData = bytes.Clone(data)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		filePath := r.PathValue(pathKey)

		// if file path not specified serve index file
		if filePath == "" || filePath == "/" {
			filePath = "index.html"
		}

		// for swagger json file
		if strings.HasSuffix(filePath, "swagger.json") {
			if docData == nil || len(docData) == 0 {
				helpers.SendError(w, http.StatusNotFound, "File not found", nil)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(docData)
			return
		}

		// for static dist files
		data, err := fs.ReadFile(fsys, path.Join(fsRootDir, filePath))
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
