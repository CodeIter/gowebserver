package handler

import (
	"fmt"
	"net/http"
)

// Redirector redirects the request to another URL.
// It extracts the domain, path and query from the request URL and constructs a new URL to redirect to.
// Redirect request /go/{domain}/{path...}?query to https://{domain}/{path...}?query
func Redirector(w http.ResponseWriter, r *http.Request) {
	proto := "https"
	domain := r.PathValue("domain")
	path := r.PathValue("path")
	query := r.URL.RawQuery

	// Construct the new URL
	newURL := fmt.Sprintf("%s://%s/%s", proto, domain, path)
	if query != "" {
		newURL += fmt.Sprintf("?%s", query)
	}

	// Perform redirection
	http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
}
