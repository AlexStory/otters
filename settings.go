package otters

type Settings struct {
	Dev       bool
	Templates TemplateSettings
	Static    StaticSettings
}

type TemplateSettings struct {
	Dir    string // directory on disk
	Layout string // filename within Dir
}

type StaticSettings struct {
	Route string
	Dir   string
}

func DefaultSettings() Settings {
	return Settings{
		Dev: true,
		Templates: TemplateSettings{
			Dir:    "./templates",
			Layout: "layout.html",
		},
		Static: StaticSettings{
			Route: "/static/",
			Dir:   "./static",
		},
	}
}
