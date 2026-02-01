package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexstory/otters"
)

func TestRecover_NonDev_Returns500AndGenericBody(t *testing.T) {
	h := Recover(false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want %q", ct, "text/plain; charset=utf-8")
	}
	if got := rr.Body.String(); got != "internal server error" {
		t.Fatalf("body = %q, want %q", got, "internal server error")
	}
}

func TestRecover_Dev_Returns500AndIncludesPanic(t *testing.T) {
	h := Recover(true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "panic: boom") {
		t.Fatalf("expected body to contain %q, got %q", "panic: boom", body)
	}
	// stack trace presence is intentionally fuzzy; just ensure there's more than the header line
	if len(body) < len("panic: boom\n")+5 {
		t.Fatalf("expected dev body to include stack trace, got %q", body)
	}

	_ = otters.Middleware(nil) // keeps import used if you later tweak tests
}
