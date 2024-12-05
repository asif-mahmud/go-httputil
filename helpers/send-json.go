package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	golog "github.com/asif-mahmud/go-log"
)

// ErrorMsg is common error message sent to the client if something went wrong
// unexpectedly while processing the request
const ErrorMsg = "Sorry, something went wrong! Please try again later."

const jsonErrorMsg = `{"status":false,"message":` + ErrorMsg + `,"data":null}`

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
