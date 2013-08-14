package main

import (
	"os/exec"
	"path"
	"strings"
	"text/template"
)

const pkgImportPath = "github.com/songgao/gallery"

func getRootPath() (string, error) {
	out, err := exec.Command("go", "list", "-f", "{{.Dir}}", "github.com/songgao/gallery").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getTemplate(filename string) (*template.Template, error) {
	root, err := getRootPath()
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(path.Join(root, "templates", filename))
}

func in(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
