package gen

import (
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

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	err = createFile("main.go", path, templates.MainTemplate, nil)
	if err != nil {
		return err
	}

	if err = createFile("go.mod", path, templates.ModTemplate, data); err != nil {
		return err
	}
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
