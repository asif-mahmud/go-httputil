package middlewares

import (
	"log/slog"
	"net/http"

	gohttputil "github.com/asif-mahmud/go-httputil"
	golog "github.com/asif-mahmud/go-log"
	"github.com/felixge/httpsnoop"
)

// LoggerWithSkips logs request and response statistics via slog for
// all routes except routes matching patterns in skipPatterns.
func LoggerWithSkips(skipPatterns ...string) gohttputil.Middleware {
	sm := map[string]any{}
	for _, s := range skipPatterns {
		sm[s] = nil
	}

	mfn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := sm[r.Pattern]; ok {
				next.ServeHTTP(w, r)
				return
			}

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
				slog.Attr{
					Key:   "pattern",
					Value: slog.StringValue(r.Pattern),
				},
			)
		}

		return http.HandlerFunc(fn)
	}

	return gohttputil.Middleware(mfn)
}
