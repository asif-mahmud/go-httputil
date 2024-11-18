package gohttputil

import "net/http"

// Middleware defines simple middleware signature
type Middleware func(http.Handler) http.Handler
