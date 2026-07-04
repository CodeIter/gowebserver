package handler

import (
	"html/template"
	"log/slog"
	assets "my-go-server"
	"net/http"
)

// templates holds the parsed HTML templates for rendering views.
var templates *template.Template

// LoadTemplates parses the HTML templates from the embedded filesystem.
func LoadTemplates() error {
	var err error
	templates, err = template.ParseFS(assets.EmbeddedFiles, "views/layout.html", "views/meta.html", "views/home.html")
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
