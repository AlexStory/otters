package otters

import (
	"net/http/httptest"
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
