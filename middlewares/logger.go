package middlewares

import (
	gohttputil "github.com/asif-mahmud/go-httputil"
)

// Logger logs request and response statistics via slog.
var Logger gohttputil.Middleware = LoggerWithSkips()
