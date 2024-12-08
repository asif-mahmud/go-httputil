package middlewares

import (
	"errors"
	"net/http"
	"strings"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/asif-mahmud/go-httputil/validator"
)

const (
	formCtxKey = "_formPayload"
	maxBytes   = 100 * 1024 * 1024
)

// ValidateForm validates request body and stores validated payload in
// the request's context.
func ValidateForm(dto any) gohttputil.Middleware {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			validatePayload(
				next,
				dto,
				func(p any) error {
					header := r.Header.Get("content-type")

					if header == "application/x-www-form-urlencoded" {
						if err := r.ParseForm(); err != nil {
							return err
						}

						return validator.BindUrlValues(r.Context(), r.Form, p)
					} else if strings.HasPrefix(header, "multipart/form-data") {
						if err := r.ParseMultipartForm(maxBytes); err != nil {
							return err
						}

						return validator.BindUrlValues(r.Context(), r.MultipartForm.Value, p)
					} else {
						return errors.New("invalid request")
					}
				},
				formCtxKey,
			).ServeHTTP(w, r)
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
