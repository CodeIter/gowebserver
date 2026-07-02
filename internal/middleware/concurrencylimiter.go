package middleware

import (
	"net/http"
)

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
