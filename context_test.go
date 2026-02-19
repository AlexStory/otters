package otters

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCtxString_WritesBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := Ctx{
		Writer:  rr,
		Request: req,
	}

	ctx.String("hello")
	if got := rr.Body.String(); got != "hello" {
		t.Errorf("String() = %v, want %v", got, "hello")
	}
}

func TestCtxRender_NoRendererReturnsError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	ctx := Ctx{
		Writer:  rr,
		Request: req,
		html:    nil,
	}

	if err := ctx.Render("index.html", nil); err == nil {
		t.Errorf("expected an error, but got nil")
	}
}

func TestCtxParam_UsesServeMuxPathValues(t *testing.T) {
	app := New()

	app.Get("/users/{id}", func(ctx *Ctx) {
		_ = ctx.String(ctx.Param("id"))
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)

	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Body.String(); got != "123" {
		t.Fatalf("body = %q, want %q", got, "123")
	}
}

func TestCtxJSON_WritesJSONAndStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := Ctx{
		Writer:  rr,
		Request: req,
	}

	type resp struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := ctx.JSON(http.StatusCreated, resp{ID: 7, Name: "otters"}); err != nil {
		t.Fatalf("JSON returned error: %v", err)
	}

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}

	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want to contain %q", ct, "application/json")
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"id":7`) || !strings.Contains(body, `"name":"otters"`) {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestCtxError_NegotiatesJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "application/json")

	ctx := Ctx{Writer: rr, Request: req}

	if err := ctx.NotFound("nope"); err != nil {
		t.Fatalf("NotFound returned error: %v", err)
	}

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want to contain %q", ct, "application/json")
	}
	if !strings.Contains(rr.Body.String(), `"error":"nope"`) {
		t.Fatalf("unexpected body: %q", rr.Body.String())
	}
}

func TestCtxError_DefaultsToHTML(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := Ctx{Writer: rr, Request: req}

	if err := ctx.NotFound("nope"); err != nil {
		t.Fatalf("NotFound returned error: %v", err)
	}

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want to contain %q", ct, "text/html")
	}
	if !strings.Contains(rr.Body.String(), "<!doctype html>") {
		t.Fatalf("expected HTML body, got: %q", rr.Body.String())
	}
}

func TestCtxJSONNotFound_ForcesJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := Ctx{Writer: rr, Request: req}

	if err := ctx.JSONNotFound("nope"); err != nil {
		t.Fatalf("JSONNotFound returned error: %v", err)
	}

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want to contain %q", ct, "application/json")
	}
}
