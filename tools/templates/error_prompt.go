package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

//go:embed error_prompt.md
var errorPrompt string

// FillErrorPrompt fills out the error prompt for the command using the text/templates
func FillErrorPrompt(err error, html, input string) (string, error) {
	errorPrompt, err := template.New("error_prompt.md").Parse(errorPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse error prompt: %w", err)
	}
	var buf bytes.Buffer
	if err := errorPrompt.Execute(&buf, struct {
		Error error
		HTML  string
		Input string
	}{
		Error: err,
		HTML:  html,
		Input: input,
	}); err != nil {
		return "", fmt.Errorf("failed to execute error prompt: %w", err)
	}
	return buf.String(), nil

}
