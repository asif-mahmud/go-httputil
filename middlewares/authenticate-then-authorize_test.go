package middlewares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/asif-mahmud/go-httputil/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateThenAuthorize(t *testing.T) {
	middlewares.SetupJWT(
		middlewares.JWTWithSecret(jwtSecret),
		middlewares.JWTWithPayloadType(user{}),
	)

	m := gohttputil.New()

	m.
		Route("/").
		Use(middlewares.Authenticate()).
		Use(middlewares.Authorize(func(r *http.Request) bool {
			u := middlewares.JWTPayload(r).(*user)
			return u.UserType == "Admin"
		})).
		Get(func(wr http.ResponseWriter, req *http.Request) {
			u := middlewares.JWTPayload(req).(*user)
			helpers.SendData(wr, u)
		})

	customerToken, _ := middlewares.DefaultJWT.Sign(
		jwt.SigningMethodHS512,
		dummyClaims(user{1, "Customer"}),
	)
	adminToken, _ := middlewares.DefaultJWT.Sign(
		jwt.SigningMethodHS512,
		dummyClaims(user{1, "Admin"}),
	)

	type testCase struct {
		token            string
		expectedStatus   int
		expectedResponse string
	}

	testCases := []testCase{
		{"", http.StatusUnauthorized, `{"data":null,"message":"Unauthorized","status":false}`},
		{
			customerToken,
			http.StatusUnauthorized,
			`{"data":null,"message":"Unauthorized","status":false}`,
		},
		{
			adminToken,
			http.StatusOK,
			`{"data":{"Id":1,"UserType":"Admin"},"message":"Success","status":true}`,
		},
	}

	for _, c := range testCases {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add("Authorization", "Bearer "+c.token)
		w := httptest.NewRecorder()

		m.ServeHTTP(w, r)

		d, e := io.ReadAll(w.Body)

		assert.Nil(t, e)
		assert.Equal(t, c.expectedStatus, w.Code)
		assert.Equal(t, c.expectedResponse, string(d))
	}
}
