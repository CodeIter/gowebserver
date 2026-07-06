package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"

	assets "github.com/CodeIter/gowebserver"
)

// templates holds the parsed HTML templates for rendering views.
var templates *template.Template

// LoadTemplates parses all HTML templates from the views directory.
func LoadTemplates() error {
	var templateFiles []string

	// Walk through the views directory and collect all .html files
	err := fs.WalkDir(assets.EmbeddedFiles, "views", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking templates directory: %w", err)
	}

	// Parse all collected template files with custom functions
	var templateErr error
	templates, templateErr = template.New("").
		Funcs(template.FuncMap{
			"renderTemplate": renderTemplateFunc,
		}).
		ParseFS(assets.EmbeddedFiles, templateFiles...)
	return templateErr
}

// RenderTemplate renders the specified template with the provided data.
func RenderTemplate(w http.ResponseWriter, tpl string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, tpl, data); err != nil {
		slog.Error("Error executing template", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// renderTemplateFunc is a template function to dynamically render templates
func renderTemplateFunc(name string, data any) (template.HTML, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		slog.Error("Error rendering template", slog.String("template", name), slog.Any("error", err))
		return "", err
	}
	return template.HTML(buf.String()), nil
}
