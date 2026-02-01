package otters

import (
	"fmt"
	"net/http"
)

type Ctx struct {
	Writer  http.ResponseWriter
	Request *http.Request

	settings Settings
	html     Renderer
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
