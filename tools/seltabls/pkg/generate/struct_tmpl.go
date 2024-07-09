package generate

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

const (
	aggregateTemplate = "aggregate"
	structTemplate    = "struct"
	promptTemplate    = "prompt"
)

// NewStructPrompt creates a new filled out template for a prompt from the `struct.tmpl` template.
func NewStructPrompt(
	url, content string,
	selectors []master.Selector,
) (string, error) {
	tmpl := template.New("struct_file_template")
	tmpl, err := tmpl.Parse(structTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse struct: %w", err)
	}
	sels := []string{}
	for _, sel := range selectors {
		sels = append(sels, fmt.Sprintf("%s: %d", sel.Value, sel.Occurances))
	}
	args := struct {
		URL       string
		Content   string
		Selectors string
		Schemas   []string
	}{
		URL:       url,
		Content:   content,
		Selectors: strings.Join(sels, "\n"),
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "struct", args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute struct file template: %w",
			err,
		)
	}
	return buf.String(), nil
}

// NewStructStruct creates a new filled out template for a struct from the `struct.tmpl` template.
func NewStructStruct(
	name, url string,
	ignoreElements []string,
	fields []Field,
) (string, error) {
	tmpl := template.New("struct_file_template")
	tmpl, err := tmpl.Parse(structTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse struct: %w", err)
	}
	args := struct {
		Name           string
		URL            string
		IgnoreElements []string
		Fields         []Field
	}{
		Name:           name,
		URL:            url,
		IgnoreElements: ignoreElements,
		Fields:         fields,
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, structTemplate, args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute struct file template: %w",
			err,
		)
	}
	return buf.String(), nil
}

// NewAggregatePrompt creates a new filled out template for a prompt from the `aggregate.tmpl` template.
func NewAggregatePrompt(
	url, content string,
	selectors []master.Selector,
	schemas []string,
) (string, error) {
	tmpl := template.New("aggregate_file_template")
	tmpl, err := tmpl.Parse(structTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse aggregate: %w", err)
	}
	sels := []string{}
	for _, sel := range selectors {
		sels = append(sels, fmt.Sprintf("%s: %d", sel.Value, sel.Occurances))
	}
	args := struct {
		URL       string
		Content   string
		Selectors string
		Schemas   []string
	}{
		URL:       url,
		Content:   content,
		Selectors: strings.Join(sels, "\n"),
		Schemas:   schemas,
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, aggregateTemplate, args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute aggregate file template: %w",
			err,
		)
	}
	return buf.String(), nil
}
