package generate

import (
	"bytes"
	"fmt"
	"text/template"

	_ "embed"
)

//go:embed identify.tmpl
var identifyTmpl string

// NewIdentifyPrompt creates a new filled out template for a prompt from the `identify.tmpl` template.
func NewIdentifyPrompt(
	url, content string,
) (string, error) {
	tmpl := template.New("identify_file_template")
	tmpl, err := tmpl.Parse(identifyTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse identify: %w", err)
	}
	// fill out the template
	args := struct {
		URL     string
		Content string
	}{
		URL:     url,
		Content: content,
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "prompt", args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute identify file template: %w",
			err,
		)
	}
	return buf.String(), nil
}

// NewIdentifyAggregatePrompt creates a new filled out template for a prompt from the `identify.tmpl` template.
func NewIdentifyAggregatePrompt(
	url, content string,
	schemas []string,
) (string, error) {
	tmpl := template.New("identify_file_template")
	tmpl, err := tmpl.Parse(identifyTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse identify: %w", err)
	}
	// fill out the template
	args := struct {
		Content string
		Schemas []string
	}{
		Schemas: schemas,
		Content: content,
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "aggregate", args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute identify file template: %w",
			err,
		)
	}
	return buf.String(), nil
}
