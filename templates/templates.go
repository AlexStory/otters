package templates

import (
	_ "embed"
)

//go:embed main.go.tmpl
var MainTemplate string

//go:embed go.mod.tmpl
var ModTemplate string

//go:embed bulma.min.css
var BulmaTemplate string
