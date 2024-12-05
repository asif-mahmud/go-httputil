package middlewares

import (
	"net/http"
	"strings"

	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret      []byte
	queryKey    string
	payloadType any
}

var DefaultJWT = &JWT{}

type JWTSetupFunc func(*JWT) *JWT

func SetupJWT(setupFuncs ...JWTSetupFunc) {
	for _, f := range setupFuncs {
		DefaultJWT = f(DefaultJWT)
	}
}

func JWTWithSecret(secret string) JWTSetupFunc {
	return func(j *JWT) *JWT {
		j.secret = []byte(strings.Clone(secret))
		return j
	}
}

func JWTWithQuery(key string) JWTSetupFunc {
	return func(j *JWT) *JWT {
		j.queryKey = strings.Clone(key)
		return j
	}
}

func JWTWithPayloadType(t any) JWTSetupFunc {
	return func(j *JWT) *JWT {
		j.payloadType = t
		return j
	}
}

const jwtPayloadKey = "_jwtPayload"

func unauthorizedResponse(w http.ResponseWriter) {
	helpers.SendError(w, http.StatusUnauthorized, "Unauthorized", nil)
}

func Authenticate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// collect token from header
		header := r.Header.Get("authorization")
		tokenStr := ""

		// collect token from query
		if len(header) == 0 && len(DefaultJWT.queryKey) > 0 {
			tokenStr = r.URL.Query().Get(DefaultJWT.queryKey)
		} else {
			tokens := strings.Split(header, " ")
			if len(tokens) != 2 {
				unauthorizedResponse(w)
				return
			}
			tokenStr = tokens[1]
		}

		// parse jwt
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(DefaultJWT.secret), nil
		})
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		if !token.Valid {
			unauthorizedResponse(w)
			return
		}

		// // set user id in the context
		// wrappedRequest := r.WithContext(context.WithValue(r.Context(), userKey, claims.Id))
		// next.ServeHTTP(w, wrappedRequest)
	}

	return http.HandlerFunc(fn)
}
