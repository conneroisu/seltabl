package templates

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/charmbracelet/huh"
	"github.com/conneroisu/seltabl"
	_ "embed"
)

//go:embed base_prompt.md
var basePrompt string

// BasePrompt fills out the base prompt for the command using the text/templates
// package.
func BasePrompt(map map[string]interface{}) (string, error) {
	
}
