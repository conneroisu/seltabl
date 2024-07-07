package prompts

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

// ErrorPrompt is the error prompt
//
//go:embed error.md
var ErrorPrompt string

// NewErrPrompt returns a new error prompt
func NewErrPrompt(html, content, url string, erro error) (string, error) {
	tmpl := template.New("errorPrompt")
	tmpl, err := tmpl.Parse(ErrorPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse error prompt: %w", err)
	}
	var data = struct {
		HTML    string
		Content string
		Err     string
		URL     string
	}{
		HTML:    html,
		Content: content,
		Err:     erro.Error(),
		URL:     url,
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute error prompt: %w", err)
	}
	return buf.String(), nil
}
