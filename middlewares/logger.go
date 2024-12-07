package middlewares

import (
	"log/slog"
	"net/http"

	golog "github.com/asif-mahmud/go-log"
	"github.com/felixge/httpsnoop"
)

// Logger logs request and response statistics via slog.
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s := httpsnoop.CaptureMetrics(next, w, r)

		slog.Info("",
			golog.Path(r.URL.Path),
			golog.Method(r.Method),
			golog.Status(s.Code),
			golog.Query(r.URL.Query()),
			golog.Ip(r.RemoteAddr),
			golog.UserAgent(r.UserAgent()),
			golog.Length(int(s.Written)),
			golog.Latency(s.Duration),
		)
	}

	return http.HandlerFunc(fn)
}
