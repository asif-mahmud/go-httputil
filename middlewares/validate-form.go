package middlewares

import (
	"net/http"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/asif-mahmud/go-httputil/validator"
)

const formCtxKey = "_formPayload"

// ValidateForm validates request body and stores validated payload in
// the request's context.
func ValidateForm(dto any) gohttputil.Middleware {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			validatePayload(
				next,
				dto,
				func(p any) error {
					return validator.BindUrlValues(r.Context(), r.Form, p)
				},
				formCtxKey,
			)
		}

		return http.HandlerFunc(fn)
	}

	return m
}

// FormPayload returns the validated form payload stored in
// request's context.
func FormPayload(r *http.Request) any {
	return r.Context().Value(formCtxKey)
}
