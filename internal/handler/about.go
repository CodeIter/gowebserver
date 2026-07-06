package handler

import (
	"net/http"

	"github.com/CodeIter/gowebserver/internal/template"
)

// About serves the about page.
func About(w http.ResponseWriter, r *http.Request) {
	template.RenderTemplate(w, "layout", map[string]string{
		"PageTemplate": "about",
		"Title":        "About Us",
		"Message":      "Learn more about our Go web server.",
		"Content":      "This is a demonstration of a Go-powered web server with modern middleware, templating, and static file serving.",
	})
}
