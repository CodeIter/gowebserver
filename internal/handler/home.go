package handler

import (
	"net/http"

	"github.com/CodeIter/gowebserver/internal/template"
)

// Home serves the home page.
func Home(w http.ResponseWriter, r *http.Request) {
	template.RenderTemplate(w, "layout", map[string]string{
		"Title":   "Welcome to Go Server",
		"Message": "Hello, World! Welcome to your Go-powered welcome page.",
	})
}
