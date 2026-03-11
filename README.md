# go-httputil

![GitHub Tag](https://img.shields.io/github/v/tag/asif-mahmud/go-httputil)
[![Go Reference](https://pkg.go.dev/badge/github.com/asif-mahmud/go-httputil.svg)](https://pkg.go.dev/github.com/asif-mahmud/go-httputil)
![GitHub License](https://img.shields.io/github/license/asif-mahmud/go-httputil)
![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/asif-mahmud/go-httputil)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/asif-mahmud/go-httputil/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/asif-mahmud/go-httputil)](https://goreportcard.com/report/github.com/asif-mahmud/go-httputil)

Simple HTTP utilities like middleware and pragmatic route builders based on standard library facilities.

This library heavily emphasizes utilization of built-in modules.

## Index

1. [Features](#features)
2. [Middleware](#middleware)
3. [Routing & Mux](#routing--mux)
4. [Validation Middlewares](#validation-middlewares)
5. [Error Formatting Design](#error-formatting-design)
6. [Authentication & Authorization](#authentication--authorization)

## Features

- Wrapper around `http.ServeMux` for global, group and route-level middleware layering.
- Bind and validate JSON, Form, Query and Path parameters.
- Built-in JWT authentication and role-based authorization middlewares.
- Automatic recovery and structured logging middlewares.
- Structural error formatting mapping go-playground/validator errors directly into nested JSON shapes.

## Middleware

In `go-httputil`, a middleware heavily embraces the standard library's definition, represented strictly as:

```go
type Middleware func(http.Handler) http.Handler
```

### Built-in Middlewares

The package includes several pragmatic middlewares out of the box (all are located under `middlewares` module):

- **`Logger()` / `LoggerWithSkips()`**: Provides structured API request logging using `log/slog`.
- **`Recover()`**: Gracefully catches panics during request handling and returns a clean 500 internal server error.
- **`Authenticate(queryKeys...)`**: Verifies JWT tokens and injects parsed payloads directly into the request context.
- **`Authorize(AuthorizeFunc)`**: Evaluates custom conditions (like RBAC) to determine if a request should proceed.
- **`Validate...()`**: A family of native validation binders for JSON, UI Forms, Queries and Path parameters.

## Routing & Mux

The `Mux` provides a thin pragmatic wrapper over Go's standard `http.ServeMux`. It allows chaining middlewares 
natively at the global, group or individual route level and handles CORS easily.

```go
mux := gohttputil.New()

// Attach global middlewares
mux.Use(middlewares.Logger(), middlewares.Recover())

// Enable default or custom CORS
mux.EnableCORS()

// Define a root route
mux.Route("/health").Get(healthCheckHandler)

// Create an authenticated sub-group
authGroup := mux.Group("/admin").Use(authMiddleware)

// Define specific nested routes using a GroupRouter closure
authGroup.Route("/users", func(r gohttputil.RouteHandler) {
    r.Use(roleMiddleware)
    r.Get(listUsersHandler)
    r.Post(createUserHandler)
})
```

## Validation Middlewares

The library provides middlewares out of the box to validate request payloads and bind them directly into contexts safely.

### Validate JSON

```go
type CreateUserDTO struct {
	Name  string `json:"name" validate:"required,min=3"`
	Age   int    `json:"age" validate:"required,gte=18"`
}

mux.Route("/users").
    Use(middlewares.ValidateJSON(CreateUserDTO{})).
    Post(func(w http.ResponseWriter, r *http.Request) {
        // Retrieve safely bound pointer from context
        payload := middlewares.JSONPayload(r).(*CreateUserDTO)
        fmt.Fprintf(w, "User %s created", payload.Name)
    })
```

### Validate Query Parameters

```go
type SearchQuery struct {
	Tags []string `form:"tags[]" validate:"required,min=1"`
}

mux.Route("/search", func(r gohttputil.RouteHandler) {
    // Bind URL queries natively
    r.Use(middlewares.ValidateQuery(SearchQuery{}))
    r.Get(handleSearch)
})
```

### Validate Form Data

```go
type SubmitForm struct {
	Title string `form:"title" validate:"required"`
}

mux.Route("/submit").
    // Bind form-data and application/x-www-form-urlencoded
    Use(middlewares.ValidateForm(SubmitForm{})).
    Post(handleSubmit)
```

> **Note:** Form validation only binds scalar values and arrays. For file uploads via `multipart/form-data` use the standard `r.FormFile()` or `r.ParseMultipartForm()` directly from the `*http.Request` in your handler.

### Validate URL Path Values

```go
type UserPathInfo struct {
	Id int `path:"id" validate:"required,gt=0"`
}

mux.Route("/users/{id}").
    Use(middlewares.ValidatePathValue(UserPathInfo{})).
    Get(func(w http.ResponseWriter, r *http.Request) {
        pathData := middlewares.PathValuePayload(r).(*UserPathInfo)
        fmt.Fprintf(w, "Fetching User ID: %d", pathData.Id)
    })
```

## Error Formatting Design

When validation fails, `go-httputil` automatically formats the validation errors to structurally
mirror the layout of your input DTO. This eliminates the need for clients to parse flat string
paths (`"Data.Addresses[0].Street"`).

### Example Structure

Given a complex nested payload:

```go
type Address struct {
	Street string `json:"street" validate:"required"`
}

type User struct {
	UserName  string    `json:"userName" validate:"required"`
	Addresses []Address `json:"addresses" validate:"required,dive"`
}
```

If the client submits an invalid payload, the middleware responds with a JSON error structure that
identically maps the error keys back to the structure:

```json
{
  "message": "Validation error",
  "status": false,
  "data": {
    "userName": "userName is a required field",
    "addresses": [
      {
        "street": "street is a required field"
      }
    ]
  }
}
```

This ensures frontend frameworks can dynamically map errors directly to form inputs using standard
JSON layouts without custom parsers.

### Extending Validation, Translations & Transformations

You can seamlessly register custom validation rules, translations, modifiers and scrubbers into the underlying
`go-playground/validator` and `go-playground/mold` engines by using the provided wrapper functions.

```go
import (
    "github.com/asif-mahmud/go-httputil/validator"
    vd "github.com/go-playground/validator/v10"
)

// 1. Register a custom validation tag
validator.RegisterValidator("is_awesome", func(fl vd.FieldLevel) bool {
    return fl.Field().String() == "awesome"
})

// 2. Register a custom error translation for a tag
validator.RegisterTranslation(validator.Translation{
    Tag:         "is_awesome",
    Translation: "{0} must be awesome",
    Override:    true,
})

// 3. Register a data modifier (e.g. data normalization before validation)
validator.RegisterModifier("uppercase", func(ctx context.Context, fl mold.FieldLevel) error {
    fl.Field().SetString(strings.ToUpper(fl.Field().String()))
    return nil
})

// 4. Register a data scrubber (e.g. sanitizing sensitive inputs)
validator.RegisterScrubber("redact", func(ctx context.Context, fl mold.FieldLevel) error {
    fl.Field().SetString("[REDACTED]")
    return nil
})
```

## Authentication & Authorization

You can configure JWT settings globally and then chain standard JWT validation alongside Role-Based
Access Control logic cleanly before your routes.

```go
type UserClaims struct {
    Role string `json:"role"`
}

// 1. Setup global JWT configuration once
middlewares.SetupJWT(
	middlewares.JWTWithSecret("secret-key"),
	middlewares.JWTWithPayloadType(UserClaims{}), // Automatically map payload into context
)

mux := gohttputil.New()

// 2. Chain Authenticate and Authorize middlewares to specific routes
mux.Route("/admin").
    // Authenticate verifies the "Authorization" Bearer token
    Use(middlewares.Authenticate()).
    // Authorize natively checks the payload to enforce custom RBAC
    Use(middlewares.Authorize(func(r *http.Request) bool {
        claims := middlewares.JWTPayload(r).(*UserClaims)
        return claims.Role == "admin"
    })).
    Get(adminDashboardHandler)
```
