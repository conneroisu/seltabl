package templates

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed base_prompt.md
var basePrompt string

// FillBasePrompt fills out the base prompt for the command using the text/templates
func FillBasePrompt(fields []string, html, input string) (string, error) {
	basePrompt, err := template.New("base_prompt.md").Parse(basePrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse base prompt: %w", err)
	}
	var buf bytes.Buffer
	if err := basePrompt.Execute(&buf, struct {
		Fields []string
		HTML   string
		Input  string
	}{
		Fields: fields,
		HTML:   html,
		Input:  input,
	}); err != nil {
		return "", fmt.Errorf("failed to execute base prompt: %w", err)
	}
	return buf.String(), nil

}
