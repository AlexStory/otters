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

//go:embed settings.go.tmpl
var SettingsTemplate string

//go:embed layout.html
var LayoutTemplate string

//go:embed index.html
var IndexTemplate string

//go:embed subapp.go.tmpl
var SubAppGoTemplate string
