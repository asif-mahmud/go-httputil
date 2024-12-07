package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/asif-mahmud/go-httputil/validator"
	golog "github.com/asif-mahmud/go-log"

	vd "github.com/go-playground/validator/v10"
)

func badrequest(w http.ResponseWriter, msg string, data any) {
	helpers.SendError(w, http.StatusBadRequest, msg, data)
}

func validatePayload(
	next http.Handler,
	dto any,
	bindFunc func(any) error,
	ctxKey string,
) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// initialize payload struct
		pv, err := helpers.NewValue(dto)
		if err != nil {
			slog.Error("Failed to initialize payload instance", golog.Extra(map[string]any{
				"error": err.Error(),
			}))
			badrequest(w, helpers.ErrorMsg, nil)
			return
		}
		p := pv.Interface()

		// bind json body
		if err := bindFunc(p); err != nil {
			slog.Error("Failed to bind payload", golog.Extra(map[string]any{
				"error": err.Error(),
			}))
			badrequest(w, helpers.ErrorMsg, nil)
			return
		}

		// validate
		if err := validator.ValidateStruct(r.Context(), p); err != nil {
			switch e := err.(type) {
			case vd.ValidationErrors:
				badrequest(w, "Validation error", validator.FormatErrors(e))
				return

			default:
				slog.Error("Failed to run validation", golog.Extra(map[string]any{
					"error": err.Error(),
				}))
				badrequest(w, helpers.ErrorMsg, nil)
				return
			}
		}

		// store in request context
		wrappedRequest := r.WithContext(context.WithValue(r.Context(), ctxKey, p))
		next.ServeHTTP(w, wrappedRequest)
	}

	return http.HandlerFunc(fn)
}
