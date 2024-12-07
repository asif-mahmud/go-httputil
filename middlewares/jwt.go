package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/asif-mahmud/go-httputil/helpers"
	golog "github.com/asif-mahmud/go-log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
)

// JWT provides interface for setting up JWT token parsing and payload retrieval.
type JWT struct {
	secret      []byte
	payloadType any
}

// Sign creates a JWT using the secret key set in setup stage.
func (j *JWT) Sign(method jwt.SigningMethod, payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, payload)
	return token.SignedString([]byte(j.secret))
}

func (j *JWT) Verify(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// DefaultJWT is global JWT instance.
var DefaultJWT = &JWT{}

// JWTSetupFunc is the signature for setting up JWT via builder function.
type JWTSetupFunc func(*JWT) *JWT

// SetupJWT sets up DefaultJWT.
// Application should setup this once in it's lifetime.
func SetupJWT(setupFuncs ...JWTSetupFunc) {
	for _, f := range setupFuncs {
		DefaultJWT = f(DefaultJWT)
	}
}

// JWTWithSecret sets JWT Secret key
func JWTWithSecret(secret string) JWTSetupFunc {
	return func(j *JWT) *JWT {
		j.secret = []byte(strings.Clone(secret))
		return j
	}
}

// JWTWithPayloadType sets the payload data type for JWT paylod.
// Upon successfull verification and parsing, request context will
// have a pointer of type t filled with data found from JWT payload.
func JWTWithPayloadType(t any) JWTSetupFunc {
	return func(j *JWT) *JWT {
		j.payloadType = t
		return j
	}
}

// jwtPayloadKey is the request context key
const jwtPayloadKey = "_jwtPayload"

// unauthorizedResponse sends unauthorized response
func unauthorizedResponse(w http.ResponseWriter) {
	helpers.SendError(w, http.StatusUnauthorized, "Unauthorized", nil)
}

// Authenticate creates a middleware to verify and parse jwt.
// By default  it will checj Bearer token from Authorization header.
// But user may specify URL search query keys in queryKeys parameter
// from which jwt can be collected, verified and parsed.
// If authentication fails an unauthorized response will be sent to
// the client.
// If authentication succeeds request's context key will contain
// JWT payload if payload type is specified in setup stage.
func Authenticate(queryKeys ...string) gohttputil.Middleware {
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// collect token from header
			header := r.Header.Get("authorization")
			tokenStr := ""

			// collect token from query if needed
			if len(header) == 0 && len(queryKeys) > 0 {
				for _, k := range queryKeys {
					if t := r.URL.Query().Get(k); len(t) > 0 {
						tokenStr = t
						break
					}
				}
			} else {
				tokens := strings.Split(header, " ")
				if len(tokens) != 2 {
					unauthorizedResponse(w)
					return
				}
				tokenStr = tokens[1]
			}

			// parse and verify jwt
			token, err := DefaultJWT.Verify(tokenStr)
			if err != nil {
				unauthorizedResponse(w)
				return
			}

			if DefaultJWT.payloadType == nil {
				next.ServeHTTP(w, r)
				return
			}

			// parse and set jwt payload in the context
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				slog.Error("Failed to parse JWT payload")
				unauthorizedResponse(w)
				return
			}

			p, err := newValue(DefaultJWT.payloadType)
			if err != nil {
				slog.Error("Failed to initiate JWT payload type", golog.Extra(map[string]any{
					"error": err.Error(),
				}))
				unauthorizedResponse(w)
				return
			}
			pi := p.Interface()
			if err := mapstructure.Decode(claims, pi); err != nil {
				slog.Error("Failed to decode payload type", golog.Extra(map[string]any{
					"error": err.Error(),
				}))
				unauthorizedResponse(w)
				return
			}
			wrappedRequest := r.WithContext(context.WithValue(r.Context(), jwtPayloadKey, pi))
			next.ServeHTTP(w, wrappedRequest)
		}

		return http.HandlerFunc(fn)
	}

	return m
}

// JWTPayload returns the context value containing the JWT payload found in authentication stage.
func JWTPayload(r *http.Request) any {
	return r.Context().Value(jwtPayloadKey)
}
