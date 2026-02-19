package otters

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Ctx struct {
	Writer  http.ResponseWriter
	Request *http.Request

	settings Settings
	html     Renderer
}

type jsonError struct {
	Error string `json:"error"`
}

// Writes the given string to the otters Context
func (c Ctx) String(content string) error {
	_, err := fmt.Fprint(c.Writer, content)
	return err
}

func (c Ctx) Render(name string, data any) error {
	if c.html == nil {
		return fmt.Errorf("html renderer is not set")
	}
	return c.html.Render(c.Writer, name, data)
}

// JSON writes the value as JSON with the given status code.
// Struct tags like `json:"name,omitempty` are respected automatically
func (c Ctx) JSON(statusCode int, v any) error {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.WriteHeader(statusCode)
	return json.NewEncoder(c.Writer).Encode(v)
}

// Param returns a path parameter extracted by Go's net/http ServeMux patterns
func (c Ctx) Param(key string) string {
	if c.Request == nil {
		return ""
	}
	return c.Request.PathValue(key)
}

// WantsJSON tries to decide whether the client expects a JSON response.
//
// uses Accept primarily Content-Type is used as a pragmatic fallback
func (c Ctx) WantsJSON() bool {
	if c.Request == nil {
		return false
	}

	accept := c.Request.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return true
	}

	ct := c.Request.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		return true
	}
	return false
}

// HTMLError renders a status-code specific error template if present (e.g. 404.html),
// otherwise it falls back to a built-in embedded HTML error page.
//
// Custom templates are looked up relative to Settings.Templates.Dir.
func (c Ctx) HTMLError(statusCode int, message string) error {
	// Try custom error template first: templates/404.html, templates/500.html, etc.
	if c.html != nil && c.settings.Templates.Dir != "" {
		name := fmt.Sprintf("%d.html", statusCode)
		fullPath := filepath.Join(c.settings.Templates.Dir, name)

		if _, err := os.Stat(fullPath); err == nil {
			// Ensure headers are set before WriteHeader.
			c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
			c.Writer.WriteHeader(statusCode)
			return c.html.Render(c.Writer, name, map[string]any{
				"StatusCode": statusCode,
				"Message":    message,
				"Request":    c.Request,
			})
		}
	}

	// Fallback to embedded default error page.
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(statusCode)

	t, err := template.New("otters_default_error").Parse(defaultErrorHTML)
	if err != nil {
		return err
	}

	return t.Execute(c.Writer, map[string]any{
		"StatusCode": statusCode,
		"Message":    message,
		"Request":    c.Request,
	})
}

// JSONError forces a JSON error envelope: {"error": "..."}.
func (c Ctx) JSONError(statusCode int, message string) error {
	return c.JSON(statusCode, jsonError{Error: message})
}

// Error chooses between JSONError and HTMLError using WantsJSON().
func (c Ctx) Error(statusCode int, message string) error {
	if c.WantsJSON() {
		return c.JSONError(statusCode, message)
	}
	return c.HTMLError(statusCode, message)
}

// Convenience helpers (negotiated)
func (c Ctx) BadRequest(message string) error   { return c.Error(http.StatusBadRequest, message) }
func (c Ctx) Unauthorized(message string) error { return c.Error(http.StatusUnauthorized, message) }
func (c Ctx) Forbidden(message string) error    { return c.Error(http.StatusForbidden, message) }
func (c Ctx) NotFound(message string) error     { return c.Error(http.StatusNotFound, message) }

// Convenience helpers (forced JSON)
func (c Ctx) JSONBadRequest(message string) error { return c.JSONError(http.StatusBadRequest, message) }
func (c Ctx) JSONUnauthorized(message string) error {
	return c.JSONError(http.StatusUnauthorized, message)
}
func (c Ctx) JSONForbidden(message string) error { return c.JSONError(http.StatusForbidden, message) }
func (c Ctx) JSONNotFound(message string) error  { return c.JSONError(http.StatusNotFound, message) }
