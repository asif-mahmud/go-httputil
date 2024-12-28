package middlewares

import (
	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/asif-mahmud/go-httputil/validator"
	"net/http"
)

const pathValueCtxKey = "_pathValuePayload"

// ValidatePathValue validates request path parameters and stores validated payload in
// the request's context.
func ValidatePathValue(dto any) gohttputil.Middleware {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			validatePayload(
				next,
				dto,
				func(p any) error {
					return validator.BindPathValues(r.Context(), r, p)
				},
				pathValueCtxKey,
			).ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return m
}

// PathValuePayload returns the validated path parameters stored in
// request's context.
func PathValuePayload(r *http.Request) any {
	return r.Context().Value(pathValueCtxKey)
}
