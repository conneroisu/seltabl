package config

import (
	"bytes"
	"fmt"
	"text/template"

	_ "embed"
)

//go:embed sections.tmpl
var sectionsTmpl string

// sectionTemplate is the template for the sections file
var sectionTemplate *template.Template

// init initializes the section template
func init() {
	var err error
	sectionTemplate = template.New("sections_file_template")
	sectionTemplate, err = sectionTemplate.Parse(sectionsTmpl)
	if err != nil {
		panic(err)
	}
}

// NewSectionsErrorPrompt creates a new filled out template for a prompt from the `sections.tmpl` template
func NewSectionsErrorPrompt(
	err error,
) (string, error) {
	// fill out the template
	args := struct {
		Error error
	}{
		Error: fmt.Errorf("failed to parse structs: %w", err),
	}
	var buf bytes.Buffer
	err = sectionTemplate.ExecuteTemplate(&buf, "error", args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute identify file template: %w",
			err,
		)
	}
	return buf.String(), nil
}
