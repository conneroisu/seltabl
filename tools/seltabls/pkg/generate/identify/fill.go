package identify

import (
	"bytes"
	"fmt"
	"text/template"

	_ "embed"
)

// identifyTmpl is the template for the identify file.
// This is the embedded template string that contains the template
// for Identity Prompt,  Aggregate Prompt, and Error Prompt.
//
//go:embed identify.tmpl
var identifyTmpl string

// identifyTemplate is the template for the identify file.
// It is used to generate identify sections in the config file.
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

// NewIdentifyPrompt creates a new filled out template for a prompt from the
// `identify.tmpl` template.
func NewIdentifyPrompt(
	url, content string,
) (result string, err error) {
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

// NewIdentifyAggregatePrompt creates a new filled out template for a prompt
// from the `identify.tmpl` template.
func NewIdentifyAggregatePrompt(
	content string,
	schemas []string,
) (result string, err error) {
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

// NewIdentifyErrorPrompt creates a new filled out template for an error prompt
// from the `identify.tmpl` template.
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
