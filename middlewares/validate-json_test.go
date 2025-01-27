package middlewares_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/asif-mahmud/go-httputil/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestValidateJSON(t *testing.T) {
	type dto struct {
		Age  float64 `json:"age"  validate:"gt=0.0"`
		Name string  `json:"name" validate:"required,min=3"`
	}

	h := middlewares.ValidateJSON(
		dto{},
	)(
		http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			d := middlewares.JSONPayload(req).(*dto)
			helpers.SendData(wr, d)
		}),
	)

	type testCase struct {
		payload           string
		expectedStatus    int
		expectedResponse  string
		contentTypeHeader string
	}

	testCases := []testCase{
		{
			`{"age":0.0,"name":""}`,
			http.StatusBadRequest,
			`{"data":{"age":"Age must be greater than 0.0","name":"Name is a required field"},"message":"Validation error","status":false}`,
			"application/json",
		},
		{
			`{"age":13.0,"name":""}`,
			http.StatusBadRequest,
			`{"data":{"name":"Name is a required field"},"message":"Validation error","status":false}`,
			"application/json",
		},
		{
			`{"age":13.5,"name":"Asif"}`,
			http.StatusOK,
			`{"data":{"age":13.5,"name":"Asif"},"message":"Success","status":true}`,
			"application/json",
		},
		{
			`{"age":13.5,"name":"Asif"}`,
			http.StatusOK,
			`{"data":{"age":13.5,"name":"Asif"},"message":"Success","status":true}`,
			"application/json; charset=UTF-8",
		},
		{
			`{"age":13.5,"name":"Asif"}`,
			http.StatusBadRequest,
			`{"data":null,"message":"Invalid request","status":false}`,
			"",
		},
	}

	for _, c := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(c.payload)))
		if len(c.contentTypeHeader) > 0 {
			r.Header.Add("Content-Type", c.contentTypeHeader)
		}
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		actual, err := io.ReadAll(w.Result().Body)

		assert.Nil(t, err)
		assert.Equal(t, c.expectedStatus, w.Result().StatusCode)
		assert.Equal(t, c.expectedResponse, string(actual))
	}
}
