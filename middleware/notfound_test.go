package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexstory/otters"
)

func TestNotFoundMiddleware_ReplacesDefaultServeMux404_WithHTML(t *testing.T) {
	app := otters.New()
	app.Use(NotFound(&app))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}

	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want to contain %q", ct, "text/html")
	}
	if !strings.Contains(rr.Body.String(), "<!doctype html>") {
		t.Fatalf("expected HTML error page, got: %q", rr.Body.String())
	}
}

func TestNotFoundMiddleware_ReplacesDefaultServeMux404_WithJSON(t *testing.T) {
	app := otters.New()
	app.Use(NotFound(&app))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	req.Header.Set("Accept", "application/json")

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}

	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want to contain %q", ct, "application/json")
	}
	if !strings.Contains(rr.Body.String(), `"error"`) {
		t.Fatalf("expected JSON error envelope, got: %q", rr.Body.String())
	}
}

func TestNotFoundMiddleware_DoesNotOverrideCustom404Body(t *testing.T) {
	app := otters.New()
	app.Use(NotFound(&app))

	app.Get("/custom", func(ctx *otters.Ctx) {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		_ = ctx.String("my custom 404")
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/custom", nil)

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	if rr.Body.String() != "my custom 404" {
		t.Fatalf("body = %q, want %q", rr.Body.String(), "my custom 404")
	}
}
