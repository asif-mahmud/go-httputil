package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/asif-mahmud/go-httputil/validator"
)

const jsonCtxKey = "_jsonPaylod"

// ValidateJSON validates JSON body and stores validated payload in
// the request's context.
func ValidateJSON(dto any) gohttputil.Middleware {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// check header conformity
			contentType := r.Header.Get("content-type")
			contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
			if contentType != "application/json" {
				slog.Error("Looking for json body, but json header is not set")
				badrequest(w, "Invalid request", nil)
				return
			}

			// bind json body
			validatePayload(
				next,
				dto,
				func(p any) error {
					return validator.BindJSON(r.Context(), r.Body, p)
				},
				jsonCtxKey,
			).ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return m
}

// JSONPayload returns the validated JSON payload stored in
// request's context.
func JSONPayload(r *http.Request) any {
	return r.Context().Value(jsonCtxKey)
}
