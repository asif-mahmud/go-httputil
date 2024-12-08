package middlewares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/asif-mahmud/go-httputil/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestValidateForm(t *testing.T) {
	type dto struct {
		Age  float64 `json:"age"  form:"age"  validate:"gt=0.0"`
		Name string  `json:"name" form:"name" validate:"required,min=3"`
	}

	h := middlewares.ValidateForm(
		dto{},
	)(
		http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			d := middlewares.FormPayload(req).(*dto)
			helpers.SendData(wr, d)
		}),
	)

	type testCase struct {
		payload          url.Values
		expectedStatus   int
		expectedResponse string
	}

	testCases := []testCase{
		{
			url.Values{"age": {"0.0"}, "name": {""}},
			http.StatusBadRequest,
			`{"data":{"age":"Age must be greater than 0.0","name":"Name is a required field"},"message":"Validation error","status":false}`,
		},
		{
			url.Values{"age": {"13.0"}, "name": {""}},
			http.StatusBadRequest,
			`{"data":{"name":"Name is a required field"},"message":"Validation error","status":false}`,
		},
		{
			url.Values{"age": {"13.5"}, "name": {"Asif"}},
			http.StatusOK,
			`{"data":{"age":13.5,"name":"Asif"},"message":"Success","status":true}`,
		},
	}

	for _, c := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(c.payload.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		actual, err := io.ReadAll(w.Result().Body)

		assert.Nil(t, err)
		assert.Equal(t, c.expectedStatus, w.Result().StatusCode)
		assert.Equal(t, c.expectedResponse, string(actual))
	}
}
