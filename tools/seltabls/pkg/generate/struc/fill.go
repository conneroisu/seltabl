package struc

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
)

const (
	aggregateTemplate = "aggregate"
	structTemplate    = "struct"
	promptTemplate    = "prompt"
)

var structTextTemplate *template.Template

// init parses the struct template upon package initialization.
func init() {
	var err error
	structTextTemplate = template.New("struct_file_template")
	structTextTemplate, err = structTextTemplate.Parse(structTmpl)
	if err != nil {
		panic(fmt.Errorf("failed to parse struct: %w", err))
	}
}

// NewStructPrompt creates a new filled out template for a prompt from the `struct.tmpl` template.
func NewStructPrompt(
	url, content string,
	selectors []master.Selector,
	section *domain.Section,
) (string, error) {
	args := struct {
		URL     string
		Content string
		Fields  []domain.Field
	}{
		URL:     url,
		Content: content,
		Fields:  section.Fields,
	}
	var buf bytes.Buffer
	err := structTextTemplate.ExecuteTemplate(&buf, "struct", args)
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
	section *domain.Section,
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
		Fields         []domain.Field
	}{
		Name:           name,
		URL:            url,
		IgnoreElements: ignoreElements,
		Fields:         section.Fields,
	}
	var buf bytes.Buffer
	err = structTextTemplate.ExecuteTemplate(&buf, structTemplate, args)
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
	err := structTextTemplate.ExecuteTemplate(&buf, aggregateTemplate, args)
	if err != nil {
		return "", fmt.Errorf(
			"failed to execute aggregate file template: %w",
			err,
		)
	}
	return buf.String(), nil
}
