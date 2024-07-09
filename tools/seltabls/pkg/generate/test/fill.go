// Package test is a package for generating a test file.
package test

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed test.tmpl
var testTemplate string

// NewTestFileContent creates a new filled out template for a test file content
func NewTestFileContent(
	name, url, version, packageName string,
) (string, error) {
	tmpl := template.New("test_file_template")
	tmpl, err := tmpl.Parse(testTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse test: %w", err)
	}
	args := struct {
		Version     string
		Name        string
		URL         string
		PackageName string
	}{
		Name:        name,
		Version:     version,
		URL:         url,
		PackageName: packageName,
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "test", args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute test file template: %w",
			err,
		)
	}
	return buf.String(), nil
}
