package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
)

// templates holds the parsed HTML templates for rendering views.
var templates *template.Template

// LoadTemplates parses the HTML templates from the views directory.
func LoadTemplates(viewsDir string) error {
	var err error
	templates, err = template.ParseFiles(
		filepath.Join(viewsDir, "layout.html"),
		filepath.Join(viewsDir, "home.html"),
	)
	return err
}

// renderTemplate renders the specified template with the provided data.
func renderTemplate(w http.ResponseWriter, tpl string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, tpl, data); err != nil {
		slog.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
