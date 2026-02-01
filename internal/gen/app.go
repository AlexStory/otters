package gen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alexstory/otters/templates"
)

type SubAppInfo struct {
	Name    string
	Package string
	Prefix  string
}

func GenApp(projectDir, name, prefix string) error {
	if name == "" {
		return fmt.Errorf("app name is required")
	}

	pkg := sanitizePackage(name)
	if prefix == "" {
		prefix = "/" + pkg
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}

	mod, err := readModuleName(projectDir)
	if err != nil {
		return err
	}

	appDir := filepath.Join(projectDir, "internal", pkg)
	if err := createDir(appDir); err != nil {
		return err
	}

	info := SubAppInfo{
		Name:    name,
		Package: pkg,
		Prefix:  prefix,
	}

	if err := createFile("app.go", appDir, templates.SubAppGoTemplate, info); err != nil {
		return err
	}

	mainPath := filepath.Join(projectDir, "main.go")
	importLine := fmt.Sprintf("\t%q", mod+"/internal/"+pkg)
	if err := insertAfterMarker(mainPath, "// ottr:imports", []string{importLine}); err != nil {
		return err
	}

	mountLine := fmt.Sprintf("\t%s.Mount(&app)", pkg)
	if err := insertAfterMarker(mainPath, "// ottr:mounts", []string{mountLine}); err != nil {
		return err
	}

	return nil
}

func sanitizePackage(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	if s == "" {
		return "app"
	}
	return s
}
