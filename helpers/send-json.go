package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	golog "github.com/asif-mahmud/go-log"
)

const jsonErrorMsg = `{"status":false,"message": "Sorry, something went wrong! Please try again later.","data":null}`

func sendJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")

	str, err := json.Marshal(data)
	if err != nil {
		slog.Error(err.Error(), golog.Extra(map[string]any{
			"data": data,
		}))

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(jsonErrorMsg))
		return
	}

	w.WriteHeader(status)
	w.Write(str)
}
