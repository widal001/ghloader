package utils

import (
	"fmt"
	"path/filepath"
	"text/template"
)

func LoadTemplate(tmplDir, tmplFile string) (*template.Template, error) {
	// Get the path to the template
	tmplPath := filepath.Join(tmplDir, tmplFile)
	// Load the template from the path
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template from %s: %v", tmplPath, err)
	}
	// Return a reference to the template
	return tmpl, nil
}
