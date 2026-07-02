package handler

import (
	"net/http"
)

// Home serves the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "layout", map[string]string{
		"Title":   "Welcome to Go Server",
		"Message": "Hello, World! Welcome to your Go-powered welcome page.",
	})
}
