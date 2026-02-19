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
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want %q", ct, "text/html; charset=utf-8")
	}

	body := rr.Body.String()
	if !strings.Contains(body, "<!doctype html>") {
		t.Fatalf("expected an HTML document, got %q", body)
	}
	if !strings.Contains(body, "panic recovered") {
		t.Fatalf("expected body to mention panic recovered, got %q", body)
	}
	if !strings.Contains(body, "boom") {
		t.Fatalf("expected body to contain %q, got %q", "boom", body)
	}

	// stack trace should usually contain "goroutine"
	if !strings.Contains(body, "goroutine") {
		t.Fatalf("expected body to include stack trace, got %q", body)
	}

	_ = otters.Middleware(nil) // keeps import used if you later tweak tests
}
