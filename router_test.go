package gohttputil_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gohttputil "github.com/asif-mahmud/go-httputil"
	"github.com/stretchr/testify/assert"
)

func handler1(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

func handler2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success":true,"method":"%s"}`, r.Method)))
}

func globalMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Global-Middleware", "1")
		h.ServeHTTP(w, r)
	})
}

func globalMiddleware1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Global-Middleware1", "1")
		h.ServeHTTP(w, r)
	})
}

func globalMiddleware2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Global-Middleware2", "2")
		h.ServeHTTP(w, r)
	})
}

func middleware1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Middleware1", "1")
		h.ServeHTTP(w, r)
	})
}

func middleware2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Middleware2", "2")
		h.ServeHTTP(w, r)
	})
}

func methodHandler(t *testing.T, method, path string) {
	m := gohttputil.New()

	m.Use(globalMiddleware)

	rh := m.Route(path).
		Use(middleware1).
		Use(middleware2)

	switch method {
	case http.MethodGet:
		rh.Get(handler1)

	case http.MethodPost:
		rh.Post(handler1)

	case http.MethodPut:
		rh.Put(handler1)

	case http.MethodPatch:
		rh.Patch(handler1)

	case http.MethodDelete:
		rh.Delete(handler1)

	default:
		assert.Fail(t, "unsupported method")
	}

	expectedBody := `{"success":true}`
	expectedHeader := http.Header{
		"Content-Type":        {"application/json"},
		"X-Global-Middleware": {"1"},
		"X-Middleware1":       {"1"},
		"X-Middleware2":       {"2"},
	}

	checkResponse(t, m, method, path, expectedBody, expectedHeader)
}

func checkResponse(
	t *testing.T,
	h http.Handler,
	method, path, expectedBody string,
	expectedHeader http.Header,
) {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	data, err := io.ReadAll(w.Body)

	assert.Nil(t, err)

	assert.Equal(t, expectedBody, string(data))

	assert.Equal(t, expectedHeader, w.Header())
}

func TestGet(t *testing.T) {
	methodHandler(t, http.MethodGet, "/")
}

func TestPost(t *testing.T) {
	methodHandler(t, http.MethodPost, "/")
}

func TestPut(t *testing.T) {
	methodHandler(t, http.MethodPut, "/")
}

func TestDelete(t *testing.T) {
	methodHandler(t, http.MethodDelete, "/")
}

func TestMultiple(t *testing.T) {
	m := gohttputil.New()

	m.Use(globalMiddleware)

	m.Route("/test").
		Get(handler2).
		Use(middleware1, middleware2).Post(handler2).
		Put(handler2).
		Patch(handler2).
		Use(middleware1, middleware2).Delete(handler2)

	type testCase struct {
		method         string
		path           string
		expectedBody   string
		expectedHeader http.Header
	}

	expectedHeader1 := http.Header{
		"Content-Type":        {"application/json"},
		"X-Global-Middleware": {"1"},
		"X-Middleware1":       {"1"},
		"X-Middleware2":       {"2"},
	}

	expectedHeader2 := http.Header{
		"Content-Type":        {"application/json"},
		"X-Global-Middleware": {"1"},
	}

	cases := []testCase{
		{http.MethodGet, "/test", `{"success":true,"method":"GET"}`, expectedHeader2},
		{http.MethodPost, "/test", `{"success":true,"method":"POST"}`, expectedHeader1},
		{http.MethodPut, "/test", `{"success":true,"method":"PUT"}`, expectedHeader2},
		{http.MethodPatch, "/test", `{"success":true,"method":"PATCH"}`, expectedHeader2},
		{http.MethodDelete, "/test", `{"success":true,"method":"DELETE"}`, expectedHeader1},
	}

	for _, c := range cases {
		checkResponse(t, m, c.method, c.path, c.expectedBody, c.expectedHeader)
	}
}
