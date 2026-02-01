package otters

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHTMLRenderer_Render_Success(t *testing.T) {
	dir := t.TempDir()

	layout := `LayoutStart|{{template "index.html" .}}|LayoutEnd`
	page := `Page:{{.Name}}`

	if err := os.WriteFile(filepath.Join(dir, "layout.html"), []byte(layout), 0o644); err != nil {
		t.Fatalf("write layout: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(page), 0o644); err != nil {
		t.Fatalf("write page: %v", err)
	}

	settings := DefaultSettings()
	settings.Templates.Dir = dir
	settings.Templates.Layout = "layout.html"

	r := NewHTMLRenderer(settings)

	rr := httptest.NewRecorder()
	err := r.Render(rr, "index.html", struct{ Name string }{Name: "Otters"})
	if err != nil {
		t.Fatalf("HTMLRenderer.Render returned error: %v", err)
	}

	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("HTMLRenderer.Render returned wrong content type: got %q want %q", ct, "text/html")
	}

	if got := rr.Body.String(); got != "LayoutStart|Page:Otters|LayoutEnd" {
		t.Errorf("HTMLRenderer.Render returned unexpected body: got %q want %q", got, "LayoutStart|Page:Otters|LayoutEnd")
	}
}

func TestHTMLRenderer_Render_EmptyTemplateDir_Error(t *testing.T) {
	settings := DefaultSettings()
	settings.Templates.Dir = ""

	r := NewHTMLRenderer(settings)

	rr := httptest.NewRecorder()
	if err := r.Render(rr, "index.html", nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestHTMLRenderer_Render_EmptyLayout_Error(t *testing.T) {
	settings := DefaultSettings()
	settings.Templates.Layout = ""

	r := NewHTMLRenderer(settings)

	rr := httptest.NewRecorder()
	if err := r.Render(rr, "index.html", nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}
