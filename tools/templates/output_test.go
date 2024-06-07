package templates

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed output_test.md
var outputTest string

// FillOutputTest fills out the output test for the command using the text/templates
func FillOutputTest(fields []string, html, input string) (string, error) {
	outputTest, err := template.New("output_test.tmpl").Parse(outputTest)
	if err != nil {
		return "", fmt.Errorf("failed to parse output test: %w", err)
	}
	var buf bytes.Buffer
	if err := outputTest.Execute(&buf, struct {
		Fields []string
		HTML   string
		Input  string
	}{
		Fields: fields,
		HTML:   html,
		Input:  input,
	}); err != nil {
		return "", fmt.Errorf("failed to execute output test: %w", err)
	}
	return buf.String(), nil

}
