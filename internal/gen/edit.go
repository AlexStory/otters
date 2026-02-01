package gen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readModuleName(projectDir string) (string, error) {
	f, err := os.Open(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("could not find module directive in go.mod")
}

func insertAfterMarker(filepath string, marker string, linesToInsert []string) error {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	content := string(b)
	idx := strings.Index(content, marker)
	if idx < 0 {
		return fmt.Errorf("marker %q not found in %s", marker, filepath)
	}

	insertPos := idx + len(marker)
	insertion := "\n" + strings.Join(linesToInsert, "\n") + "\n"
	out := content[:insertPos] + insertion + content[insertPos:]
	return os.WriteFile(filepath, []byte(out), 0644)
}
