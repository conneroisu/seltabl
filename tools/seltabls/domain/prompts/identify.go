package prompts

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed identify.md
var IdentifyPrompt string

// NewIdentifyPrompt returns a new identify prompt
// This prompt is used to identify the information that can be extracted from two given urls.
func NewIdentifyPrompt(url1, url2, url1Content, url2Content string) (string, error) {
	tmpl := template.New("identifyPrompt")
	tmpl, err := tmpl.Parse(IdentifyPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse identify prompt: %w", err)
	}
	var data = struct {
		URL1        string
		URL2        string
		URL1Content string
		URL2Content string
	}{
		URL1:        url1,
		URL2:        url2,
		URL1Content: url1Content,
		URL2Content: url2Content,
	}
	w := new(bytes.Buffer)
	err = tmpl.Execute(w, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute identify prompt: %w", err)
	}
	return w.String(), nil
}
