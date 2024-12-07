package middlewares

import (
	"net/http"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/asif-mahmud/go-httputil/validator"
)

const queryCtxKey = "_queryPayload"

// ValidateQuery validates request search query and stores validated payload in
// the request's context.
func ValidateQuery(dto any) gohttputil.Middleware {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			validatePayload(
				next,
				dto,
				func(p any) error {
					return validator.BindUrlValues(r.Context(), r.URL.Query(), p)
				},
				queryCtxKey,
			)
		}

		return http.HandlerFunc(fn)
	}

	return m
}

// QueryPayload returns the validated form payload stored in
// request's context.
func QueryPayload(r *http.Request) any {
	return r.Context().Value(queryCtxKey)
}
