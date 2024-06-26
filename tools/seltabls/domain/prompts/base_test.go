package prompts

import "testing"

// TestNewBasePrompt tests the NewBasePrompt function
func TestNewBasePrompt(t *testing.T) {
	selectors := []string{
		"html",
		"html head",
		"html body",
		"html body table",
		"html body table tbody",
		"html body table tbody tr",
		"html body table tbody tr td",
	}
	content := "Hello, World!"
	url := "https://example.com"
	prompt, err := NewBasePrompt(selectors, content, url)
	if err != nil {
		t.Errorf("NewBasePrompt() error = %v", err)
		t.Fail()
	}
	t.Logf("prompt: %s", prompt)
}
