package otters

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalizeMountPrefix(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", "/"},
		{"/", "/"},
		{"admin", "/admin"},
		{"/admin", "/admin"},
		{"/admin/", "/admin"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := normalizeMountPrefix(tt.in)
			if got != tt.want {
				t.Errorf("normalizedMountPrefix(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestAppGet_RunsHandler(t *testing.T) {
	app := New()

	app.Get("/ping", func(ctx *Ctx) {
		ctx.String("ok")
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Body.String(); got != "ok" {
		t.Fatalf("body = %q, want %q", got, "ok")
	}
}
