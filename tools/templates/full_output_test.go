package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

//go:embed full_output_test.tmpl
var fullOutputTest string

// FillFullOutputTest fills out the full output test for the command using the text/templates
func FillFullOutputTest(fields []string, html, input string) (string, error) {
	fullOutputTest, err := template.New("full_output_test.tmpl").Parse(fullOutputTest)
	if err != nil {
		return "", fmt.Errorf("failed to parse full output test: %w", err)
	}
	var buf bytes.Buffer
	if err := fullOutputTest.Execute(&buf, struct {
		Fields []string
		HTML   string
		Input  string
	}{
		Fields: fields,
		HTML:   html,
		Input:  input,
	}); err != nil {
		return "", fmt.Errorf("failed to execute full output test: %w", err)
	}
	return buf.String(), nil

}
