package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger is a middleware that logs the start and end of each request
// along with its duration.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("request started", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
		slog.Info("request completed", "duration", time.Since(start))
	})
}

// SecurityHeaders adds common security headers to the response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

// ConcurrencyLimiter uses a buffered channel as a semaphore
// to limit the number of concurrent requests being processed.
func ConcurrencyLimiter(max int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		sem := make(chan struct{}, max)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
				next.ServeHTTP(w, r)
			default:
				http.Error(w, `{"error":"Service Unavailable","message":"Server overloaded","code":503}`, http.StatusServiceUnavailable)
			}
		})
	}
}
