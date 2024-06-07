package templates

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed output.tmpl
var output string

// FillOutput fills out the output for the command using the text/templates
func FillOutput(fields []string, html, input string) (string, error) {
	output, err := template.New("output.tmpl").Parse(output)
	if err != nil {
		return "", fmt.Errorf("failed to parse output: %w", err)
	}
	var buf bytes.Buffer
	if err := output.Execute(&buf, struct {
		Fields []string
		HTML   string
		Input  string
	}{
		Fields: fields,
		HTML:   html,
		Input:  input,
	}); err != nil {
		return "", fmt.Errorf("failed to execute output: %w", err)
	}
	return buf.String(), nil

}
