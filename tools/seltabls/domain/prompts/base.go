package prompts

import (
	"fmt"
	"strings"

	_ "embed"
	"text/template"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

// BasePrompt is the base prompt for the LLM
//
//go:embed base.md
var BasePrompt string

// NewBasePrompt returns a new base prompt
func NewBasePrompt(
	selectors []master.Selector,
	content, url string,
) (string, error) {
	t, err := template.New("base").Parse(BasePrompt)
	if err != nil {
		return "", fmt.Errorf("failed to parse base prompt: %w", err)
	}
	type Args struct {
		Selectors string
		Content   string
		URL       string
	}
	sels := []string{}
	for _, sel := range selectors {
		sels = append(sels, sel.Value)
	}
	args := Args{
		Selectors: strings.Join(sels, "\n"),
		Content:   content,
		URL:       url,
	}
	buf := strings.Builder{}
	err = t.Execute(&buf, args)
	if err != nil {
		return "", fmt.Errorf("failed to execute base prompt: %w", err)
	}
	return buf.String(), nil
}
