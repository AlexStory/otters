package otters

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any) error
}

type HTMLRenderer struct {
	settings Settings
}

func NewHTMLRenderer(settings Settings) *HTMLRenderer {
	return &HTMLRenderer{settings: settings}
}

func (r *HTMLRenderer) Render(w http.ResponseWriter, name string, data any) error {
	templateDir := r.settings.Templates.Dir
	layout := r.settings.Templates.Layout

	if templateDir == "" {
		return fmt.Errorf("template dir is empty")
	}

	if layout == "" {
		return fmt.Errorf("layout template is not set")
	}

	layoutPath := filepath.Join(templateDir, layout)
	pagePath := filepath.Join(templateDir, name)

	t, err := template.ParseFiles(layoutPath, pagePath)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteTemplate(w, layout, data)
}
