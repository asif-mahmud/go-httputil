package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/asif-mahmud/go-httputil/helpers"
	golog "github.com/asif-mahmud/go-log"
)

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error(
					fmt.Sprintf("Recovered from panic. err: %v", err),
					golog.Extra(map[string]any{
						"stack": string(debug.Stack()),
					}),
				)
				helpers.SendError(w, http.StatusBadRequest, helpers.ErrorMsg, nil)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
