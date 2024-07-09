package generate

import (
	"bytes"
	"fmt"
	"text/template"

	_ "embed"
)

//go:embed identify.tmpl
var identifyTmpl string
var identifyTemplate *template.Template

// init initializes the identify template.
func init() {
	var err error
	identifyTemplate = template.New("identify_file_template")
	identifyTemplate, err = identifyTemplate.Parse(identifyTmpl)
	if err != nil {
		panic(err)
	}
}

// NewIdentifyPrompt creates a new filled out template for a prompt from the `identify.tmpl` template.
func NewIdentifyPrompt(
	url, content string,
) (result string, err error) {
	// fill out the template
	args := struct {
		URL     string
		Content string
	}{
		URL:     url,
		Content: content,
	}
	var buf bytes.Buffer
	err = identifyTemplate.ExecuteTemplate(&buf, "prompt", args)
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
	content string,
	schemas []string,
) (result string, err error) {
	// fill out the template
	args := struct {
		Content string
		Schemas []string
	}{
		Schemas: schemas,
		Content: content,
	}
	var buf bytes.Buffer
	err = identifyTemplate.ExecuteTemplate(&buf, "aggregate", args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute identify file template: %w",
			err,
		)
	}
	return buf.String(), nil
}

// NewIdentifyErrorPrompt creates a new filled out template for an error prompt from the `identify.tmpl` template.
func NewIdentifyErrorPrompt(
	err error,
) (string, error) {
	args := struct {
		Err error
	}{
		Err: err,
	}
	var buf bytes.Buffer
	err = identifyTemplate.ExecuteTemplate(&buf, "error", args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute identify file template: %w",
			err,
		)
	}
	return buf.String(), nil
}
