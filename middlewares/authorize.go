package middlewares

import (
	"net/http"

	gohttputil "github.com/asif-mahmud/go-httputil"
)

// AuthorizeFunc is the signature for request authorizing function.
type AuthorizeFunc func(*http.Request) bool

// Authorize checks if the request is permitted to reach the handler.
// If AuthorizeFunc f returns true the request reaches the handler
// else it responds unauthorized to the client.
func Authorize(f AuthorizeFunc) gohttputil.Middleware {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !f(r) {
				unauthorizedResponse(w)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return m
}
