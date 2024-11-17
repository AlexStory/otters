package otters

import (
	"fmt"
	"net/http"
)

type Ctx struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

// Writes the given string to the otters Context
func (c Ctx) String(content string) {
	fmt.Fprint(c.Writer, content)
}
