package gen

import (
	"os"
	"path/filepath"
	"text/template"
)

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
