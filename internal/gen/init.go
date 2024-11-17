package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

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

func createFile(name, path, templateString string, data any) error {
	t, err := template.New("temp").Parse(templateString)
	if err != nil {
		return err
	}

	filename := filepath.Join(path, name)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	return t.Execute(f, data)
}

func createDir(path string) error {
	return os.MkdirAll(path, 0755)
}
