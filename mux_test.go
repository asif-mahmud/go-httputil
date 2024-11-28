package gohttputil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
)

func TestMuxCreate(t *testing.T) {
	m := gohttputil.New()

	assert.NotNil(t, m)
}

func TestAllowAllCORS(t *testing.T) {
	m := gohttputil.New()

	m.EnableCORS()

	m.Route("/").Get(handler1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodOptions, "/", nil)

	r.Header.Add("Host", "foo.bar")
	r.Header.Add("Access-Control-Request-Method", "GET")
	r.Header.Add("Access-Control-Request-Headers", "X-Header")
	r.Header.Add("Origin", "https://foo.bar")

	m.ServeHTTP(w, r)

	expectedStatusCode := http.StatusNoContent
	expected := http.Header{
		"Access-Control-Allow-Headers": {"X-Header"},
		"Access-Control-Allow-Methods": {http.MethodGet},
		"Access-Control-Allow-Origin":  {"*"},
		"Vary": {
			"Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		},
	}

	assert.Equal(t, expectedStatusCode, http.StatusNoContent)
	assert.Equal(t, expected, w.Header())
}

func TestRestrictiveCORS(t *testing.T) {
	m := gohttputil.New()

	m.Route("/").Get(handler1)

	m.EnableCORS(cors.Options{
		AllowedHeaders: []string{"X-Snoop"},
		AllowedOrigins: []string{"none.allowed.com"},
		AllowedMethods: []string{http.MethodPut},
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodOptions, "/", nil)

	r.Header.Add("Access-Control-Request-Method", "PUT")
	r.Header.Add("Access-Control-Request-Headers", "X-Header")
	r.Header.Add("Origin", "https://foo.bar")

	m.ServeHTTP(w, r)

	expectedStatusCode := http.StatusNoContent
	expected := http.Header{
		"Vary": {
			"Origin, Access-Control-Request-Method, Access-Control-Request-Headers",
		},
	}

	assert.Equal(t, expectedStatusCode, http.StatusNoContent)
	assert.Equal(t, expected, w.Header())
}
