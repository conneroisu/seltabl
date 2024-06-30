package prompts

import (
	"testing"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

// TestNewBasePrompt tests the NewBasePrompt function
func TestNewBasePrompt(t *testing.T) {
	selectors := []master.Selector{
		{
			Value:   "html",
			UrlID:   0,
			Context: "html",
		},
		{
			Value:   "html head",
			UrlID:   0,
			Context: "html head",
		},
		{
			Value:   "html body",
			UrlID:   0,
			Context: "html body",
		},
		{
			Value:   "html body table",
			UrlID:   0,
			Context: "html body table",
		},
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
