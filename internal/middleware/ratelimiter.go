package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"my-go-server/pkg/response"
)

// IPRateLimiter manages rate limiters for individual IPs.
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit // Requests per second
	b   int        // Burst size (max immediate requests)
}

// NewIPRateLimiter creates a new per-IP rate limiter.
// r: average rate (requests per second)
// b: burst capacity (max tokens in bucket)
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	rl := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
	// Start cleanup goroutine to prevent memory leaks
	go rl.cleanup()
	return rl
}

// getLimiter returns the rate limiter for a specific IP.
// It creates a new one if it doesn't exist.
func (rl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.ips[ip] = limiter
	}
	return limiter
}

// Allow checks if a request is allowed for the given IP.
func (rl *IPRateLimiter) Allow(ip string) bool {
	return rl.getLimiter(ip).Allow()
}

// cleanup removes stale IP entries older than 3 minutes.
// This prevents unbounded memory growth.
func (rl *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, limiter := range rl.ips {
			// If the limiter hasn't been used recently (tokens are full), remove it.
			// Note: rate.Limiter doesn't expose lastUsed, so we rely on token count
			// or track lastSeen separately in a production map if strictness is needed.
			// For simplicity, we check if tokens are full (idle).
			if limiter.Tokens() == float64(rl.b) {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// extractIP gets the real client IP, handling proxied requests.
// TODO move to pkg/clientip/ dir
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For header first (set by load balancers/proxies)
	xf := r.Header.Get("X-Forwarded-For")
	if xf != "" {
		// Take the first IP in the list (original client)
		parts := strings.Split(xf, ",")
		return strings.TrimSpace(parts[0])
	}

	// Check X-Real-IP (common in Nginx setups)
	xr := r.Header.Get("X-Real-IP")
	if xr != "" {
		return xr
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	// Strip port number if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		return ip[:idx]
	}
	return ip
}

// RateLimiterMiddleware creates the HTTP middleware.
func RateLimiterMiddleware(limit int, burstCapacity int) func(http.Handler) http.Handler {
	// Convert int to rate.Limit (float64)
	rl := NewIPRateLimiter(rate.Limit(limit), burstCapacity)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r)

			if !rl.Allow(ip) {
				// Set standard rate limit headers
				w.Header().Set("X-RateLimit-Limit", string(rune(limit)))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "1")

				response.Error(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}

			// Optional: Inform client of remaining quota
			limiter := rl.getLimiter(ip)
			w.Header().Set("X-RateLimit-Limit", string(rune(limit)))
			w.Header().Set("X-RateLimit-Remaining", string(rune(int(limiter.Tokens()))))

			next.ServeHTTP(w, r)
		})
	}
}
