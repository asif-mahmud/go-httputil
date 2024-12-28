package middlewares_test

import (
	"fmt"
	"github.com/asif-mahmud/go-httputil/helpers"
	_ "github.com/asif-mahmud/go-httputil/helpers"
	"github.com/asif-mahmud/go-httputil/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"
)

type dto struct {
	Id   int    `json:"id" path:"id" validate:"gt=0"`
	Name string `json:"name" path:"name" validate:"required,min=3"`
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	params := middlewares.PathValuePayload(r).(*dto)
	helpers.SendData(w, params)
}

func performRequest(r http.Handler, method, url string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestValidatePathMiddleware(t *testing.T) {
	validURL := "/users/123/John"

	// Initialize the router
	mux := http.NewServeMux()

	// Add the route with validation middleware
	mux.Handle("/users/{id}/{name}", middlewares.ValidatePathValue(
		dto{},
	)(http.HandlerFunc(getUserHandler)))

	// Perform the valid request
	rr := performRequest(mux, "GET", validURL)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code 200, got %v", status)
	}

	expectedResponse := `{"data":{"id":123,"name":"John"},"message":"Success","status":true}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %v, got %v", expectedResponse, rr.Body.String())
	}
}

// Test ValidatePath middleware with invalid name
func TestValidatePathMiddlewareInvalidName(t *testing.T) {
	// Define an invalid URL path (name too short)
	invalidURL := "/users/123/Jo"

	// Initialize the router
	mux := http.NewServeMux()

	// Add the route with validation middleware
	mux.Handle("/users/{id}/{name}", middlewares.ValidatePathValue(
		dto{},
	)(http.HandlerFunc(getUserHandler)))

	// Perform the invalid request
	rr := performRequest(mux, "GET", invalidURL)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %v", status)
	}

	expectedResponse := `{"data":{"name":"Name must be at least 3 characters in length"},"message":"Validation error","status":false}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %v, got %v", expectedResponse, rr.Body.String())
	}
}
