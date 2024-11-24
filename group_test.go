package gohttputil_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/stretchr/testify/assert"
)

func groupMiddleware1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Group-Middleware1", "1")
		h.ServeHTTP(w, r)
	})
}

func groupMiddleware2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Group-Middleware2", "2")
		h.ServeHTTP(w, r)
	})
}

func TestGroup(t *testing.T) {
	m := gohttputil.New()

	m.Use(globalMiddleware)

	m.Group("/api/v1").
		Route("/1", func(rh gohttputil.RouteHandler) {
			rh.
				Use(middleware1).
				Get(handler2).
				Use(middleware1, middleware2).
				Post(handler2)
		}).
		Use(groupMiddleware1).
		Route("/2", func(rh gohttputil.RouteHandler) {
			rh.
				Use(middleware1).
				Get(handler2).
				Use(middleware1, middleware2).
				Post(handler2)
		}).
		Use(groupMiddleware2).
		Route("/3", func(rh gohttputil.RouteHandler) {
			rh.
				Use(middleware1).
				Get(handler2).
				Use(middleware1, middleware2).
				Post(handler2)
		})

	type testCase struct {
		route          string
		method         string
		expectedResp   string
		expectedHeader http.Header
	}

	cases := []testCase{
		{"/api/v1/1", http.MethodGet, `{"success":true,"method":"GET"}`, http.Header{
			"Content-Type":        {"application/json"},
			"X-Global-Middleware": {"1"},
			"X-Middleware1":       {"1"},
		}},

		{"/api/v1/1", http.MethodPost, `{"success":true,"method":"POST"}`, http.Header{
			"Content-Type":        {"application/json"},
			"X-Global-Middleware": {"1"},
			"X-Middleware1":       {"1"},
			"X-Middleware2":       {"2"},
		}},

		{"/api/v1/2", http.MethodGet, `{"success":true,"method":"GET"}`, http.Header{
			"Content-Type":        {"application/json"},
			"X-Global-Middleware": {"1"},
			"X-Group-Middleware1": {"1"},
			"X-Middleware1":       {"1"},
		}},

		{"/api/v1/2", http.MethodPost, `{"success":true,"method":"POST"}`, http.Header{
			"Content-Type":        {"application/json"},
			"X-Global-Middleware": {"1"},
			"X-Group-Middleware1": {"1"},
			"X-Middleware1":       {"1"},
			"X-Middleware2":       {"2"},
		}},

		{"/api/v1/3", http.MethodGet, `{"success":true,"method":"GET"}`, http.Header{
			"Content-Type":        {"application/json"},
			"X-Global-Middleware": {"1"},
			"X-Group-Middleware1": {"1"},
			"X-Group-Middleware2": {"2"},
			"X-Middleware1":       {"1"},
		}},

		{"/api/v1/3", http.MethodPost, `{"success":true,"method":"POST"}`, http.Header{
			"Content-Type":        {"application/json"},
			"X-Global-Middleware": {"1"},
			"X-Group-Middleware1": {"1"},
			"X-Group-Middleware2": {"2"},
			"X-Middleware1":       {"1"},
			"X-Middleware2":       {"2"},
		}},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(c.method, c.route, nil)

		http.DefaultServeMux.ServeHTTP(w, r)

		data, err := io.ReadAll(w.Body)

		assert.Nil(t, err)

		assert.Equal(t, c.expectedResp, string(data))

		assert.Equal(t, c.expectedHeader, w.Header())
	}
}
