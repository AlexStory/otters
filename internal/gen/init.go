package gen

import (
	"fmt"
	"path/filepath"

	"github.com/alexstory/otters/templates"
)

type MainInfo struct {
	Name string
}

func Init(name, path string) error {
	data := MainInfo{
		Name: name,
	}

	fmt.Printf("creating %s...\n", path)
	if err := createDir(path); err != nil {
		return err
	}

	staticDir := filepath.Join(path, "static")
	fmt.Printf("creating %s...\n", staticDir)
	if err := createDir(staticDir); err != nil {
		return err
	}

	templateDir := filepath.Join(path, "templates")
	fmt.Printf("creating %s...\n", templateDir)
	if err := createDir(templateDir); err != nil {
		return err
	}

	fmt.Printf("creating %s...\n", filepath.Join(templateDir, "layout.html"))
	if err := createRawFile("layout.html", templateDir, templates.LayoutTemplate); err != nil {
		return err
	}

	fmt.Printf("creating %s...\n", filepath.Join(templateDir, "index.html"))
	if err := createRawFile("index.html", templateDir, templates.IndexTemplate); err != nil {
		return err
	}

	fmt.Printf("creating %s...\n", filepath.Join(staticDir, "bulma.min.css"))
	if err := createRawFile("bulma.min.css", staticDir, templates.BulmaTemplate); err != nil {
		return err
	}

	fmt.Printf("creating %s...\n", filepath.Join(path, "settings.go"))
	if err := createFile("settings.go", path, templates.SettingsTemplate, nil); err != nil {
		return err
	}

	fmt.Printf("creating %s...\n", filepath.Join(path, "main.go"))
	err := createFile("main.go", path, templates.MainTemplate, nil)
	if err != nil {
		return err
	}

	fmt.Printf("creating %s...\n", filepath.Join(path, "go.mod"))
	if err = createFile("go.mod", path, templates.ModTemplate, data); err != nil {
		return err
	}

	fmt.Println("all files created successfully!")
	fmt.Println("run the following:")
	fmt.Printf("cd %s\n", path)
	fmt.Println("go mod tidy")
	return nil
}
