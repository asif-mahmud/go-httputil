package middlewares_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/asif-mahmud/go-httputil/middlewares"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const (
	jwtSecret       = "123456790Secrit"
	invalidTokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

type myClaims struct {
	Id       int
	UserType string
	jwt.RegisteredClaims
}

func dummyClaims() myClaims {
	n := time.Now()
	return myClaims{
		Id:       1,
		UserType: "Customer",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "1",
			Subject:   "1",
			ExpiresAt: jwt.NewNumericDate(n.Add(1 * time.Second)),
			NotBefore: jwt.NewNumericDate(n),
			IssuedAt:  jwt.NewNumericDate(n),
			ID:        "1",
		},
	}
}

func TestSigning(t *testing.T) {
	middlewares.SetupJWT(middlewares.JWTWithSecret(jwtSecret))

	expected, err := middlewares.DefaultJWT.Sign(jwt.SigningMethodHS512, dummyClaims())

	assert.Nil(t, err)
	assert.NotZero(t, expected)
}

func TestVerifySuccess(t *testing.T) {
	middlewares.SetupJWT(middlewares.JWTWithSecret(jwtSecret))

	expected, _ := middlewares.DefaultJWT.Sign(jwt.SigningMethodHS512, dummyClaims())

	token, err := middlewares.DefaultJWT.Verify(expected)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestVerifyFailed(t *testing.T) {
	middlewares.SetupJWT(middlewares.JWTWithSecret(jwtSecret))

	expected, _ := middlewares.DefaultJWT.Sign(jwt.SigningMethodHS512, dummyClaims())

	// exceed expiry
	time.Sleep(1 * time.Second)

	token, err := middlewares.DefaultJWT.Verify(expected)

	assert.NotNil(t, err)
	assert.Nil(t, token)

	token, err = middlewares.DefaultJWT.Verify(invalidTokenStr)

	assert.NotNil(t, err)
	assert.Nil(t, token)
}

func TestNoPaylodTypeSuccess(t *testing.T) {
	middlewares.SetupJWT(middlewares.JWTWithSecret(jwtSecret))

	expected, _ := middlewares.DefaultJWT.Sign(jwt.SigningMethodHS512, dummyClaims())

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", expected))

	w := httptest.NewRecorder()

	m := middlewares.Authenticate()

	m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

func TestNoPaylodTypeFailed(t *testing.T) {
	middlewares.SetupJWT(middlewares.JWTWithSecret(jwtSecret))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", invalidTokenStr))

	w := httptest.NewRecorder()

	m := middlewares.Authenticate()

	m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestWithPaylodType(t *testing.T) {
	type user struct {
		Id       int
		UserType string
	}
	middlewares.SetupJWT(
		middlewares.JWTWithSecret(jwtSecret),
		middlewares.JWTWithPayloadType(user{}),
	)

	token, _ := middlewares.DefaultJWT.Sign(jwt.SigningMethodHS512, dummyClaims())

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	w := httptest.NewRecorder()

	m := middlewares.Authenticate()

	// response body will contain payload found from jwt
	m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.SendData(w, middlewares.JWTPayload(r))
	})).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	expected := `{"data":{"Id":1,"UserType":"Customer"},"message":"Success","status":true}`
	actual, err := io.ReadAll(w.Body)

	assert.Nil(t, err)
	assert.Equal(t, expected, string(actual))
}
