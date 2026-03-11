package middlewares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/asif-mahmud/go-httputil/helpers"
	"github.com/asif-mahmud/go-httputil/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestValidateQuery(t *testing.T) {
	type dto struct {
		Age  float64 `json:"age"  form:"age"  validate:"gt=0.0"`
		Name string  `json:"name" form:"name" validate:"required,min=3"`
	}

	h := middlewares.ValidateQuery(
		dto{},
	)(
		http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			d := middlewares.QueryPayload(req).(*dto)
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
			`{"data":{"age":"age must be greater than 0.0","name":"name is a required field"},"message":"Validation error","status":false}`,
		},
		{
			url.Values{"age": {"13.0"}, "name": {""}},
			http.StatusBadRequest,
			`{"data":{"name":"name is a required field"},"message":"Validation error","status":false}`,
		},
		{
			url.Values{"age": {"13.5"}, "name": {"Asif"}},
			http.StatusOK,
			`{"data":{"age":13.5,"name":"Asif"},"message":"Success","status":true}`,
		},
	}

	for _, c := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/?"+c.payload.Encode(), nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, r)

		actual, err := io.ReadAll(w.Result().Body)

		assert.Nil(t, err)
		assert.Equal(t, c.expectedStatus, w.Result().StatusCode)
		assert.Equal(t, c.expectedResponse, string(actual))
	}
}

func TestValidateQuery_Slice(t *testing.T) {
	type DTO struct {
		Values []string `form:"value[]" validate:"required,min=2"`
	}

	h := middlewares.ValidateQuery(DTO{})(
		http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			d := middlewares.QueryPayload(req).(*DTO)
			helpers.SendData(wr, d)
		}),
	)

	// Build query with multiple values
	q := url.Values{}
	q.Add("value[]", "test1")
	q.Add("value[]", "test2")

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}
