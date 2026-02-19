package otters

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

	root := os.DirFS(templateDir)
	layoutPath := filepath.ToSlash(layout)
	pagePath := filepath.ToSlash(name)
	var partials []string

	if err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, "_") {
			partials = append(partials, filepath.ToSlash(path))
		}
		return nil
	}); err != nil {
		return err
	}

	paths := make([]string, 0, 2+len(partials))
	paths = append(paths, layoutPath)

	for _, p := range partials {
		if p == layoutPath || p == pagePath {
			continue
		}
		paths = append(paths, p)
	}
	paths = append(paths, pagePath)

	t := template.New(filepath.Base(layout))
	t, err := t.ParseFS(root, paths...)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteTemplate(w, layout, data)
}
