package helpers

import (
	"encoding/json"
	"log/slog"

	golog "github.com/asif-mahmud/go-log"
)

// ToJson marshals any data into json string and ignores any error.
// On error, it logs the error and returns empty string.
func ToJson(data any) string {
	j, e := json.Marshal(data)
	if e != nil {
		slog.Error("error marshaling data to json", golog.Extra(map[string]any{
			"error": e.Error(),
		}))
		return ""
	}
	return string(j)
}
