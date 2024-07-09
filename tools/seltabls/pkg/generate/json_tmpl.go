package generate

import (
	"bytes"
	"fmt"
	"text/template"

	_ "embed"
)

//go:embed json.tmpl
var jsonTmpl string

// NewJSONPrompt creates a new filled out template for a prompt from the `json.tmpl` template.
func NewJSONPrompt(content string) (string, error) {
	tmpl := template.New("json_file_template")
	tmpl, err := tmpl.Parse(jsonTmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse json: %w", err)
	}
	args := struct {
		Content string
	}{
		Content: content,
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "prompt", args)
	if err != nil {
		return "", fmt.Errorf("failed to execute json file template: %w", err)
	}
	return buf.String(), nil
}
