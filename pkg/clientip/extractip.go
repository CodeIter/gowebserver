package clientip

import (
	"net/http"
	"strings"
)

// ExtractIP gets the real client IP, handling proxied requests.
func ExtractIP(r *http.Request) string {
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
