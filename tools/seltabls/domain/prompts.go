package domain

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/liushuangls/go-anthropic/v2"
)

//go:embed prompts.tmpl
var promptsTmpl string

// PromptTemplate is the must compile template
var PromptTemplate = template.Must(template.New("templates").Parse(promptsTmpl))

// prompter is an interface for prompting it provided the template name to use.
type prompter interface {
	prompt() string
}

// responder is an interface for responding to a given prompt.
type responder interface{ respond() string }

func (i IdentifyResponse) respond() string {
	return "identify_response"
}

// NewPrompt creates a new prompt for the given args
func NewPrompt(
	promptArgs prompter,
) (string, error) {
	name := promptArgs.prompt()
	buf := new(bytes.Buffer)
	err := PromptTemplate.ExecuteTemplate(buf, name, promptArgs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// PromptBetterError is the arguments for the better error prompt.
type PromptBetterError struct {
	Error   error
	History []anthropic.Message
}

func (a PromptBetterError) prompt() string {
	return "section_error"
}

// PromptAggregateSections is the arguments for the section aggregate prompt.
type PromptAggregateSections struct {
	Structs   []string          `json:"structs"`   // required
	Content   string            `json:"content"`   // required
	Selectors []master.Selector `json:"selectors"` // required
}

func (a PromptAggregateSections) prompt() string { return "section_aggregate" }

// PromptAggregateStructs is the arguments for the struct aggregate prompt.
type PromptAggregateStructs struct {
	Selectors []master.Selector `json:"selectors"` // required
	Content   string            `json:"content"`   // required
	Schemas   []string          `json:"schemas"`   // required
}

func (a PromptAggregateStructs) prompt() string { return "struct_aggregate" }

// IdentifyArgs is the arguments for the identify prompt.
type IdentifyArgs struct {
	URL         string            `json:"url"`
	Content     string            `json:"content"`
	NumSections int               `json:"num-sections"`
	Selectors   []master.Selector `json:"selectors"`
}

func (a IdentifyArgs) prompt() string { return "identify_prompt" }

// IdentifyAggregateArgs is the arguments for the identify aggregate prompt.
type IdentifyAggregateArgs struct {
	Content   string            `json:"content"`
	Schemas   []string          `json:"schemas"`
	Selectors []master.Selector `json:"selectors"`
}

func (a IdentifyAggregateArgs) prompt() string { return "identify_aggregate" }

// StructPromptArgs is the arguments for the struct prompt.
type StructPromptArgs struct {
	URL       string            `json:"url"`
	Content   string            `json:"content"`
	Selectors []master.Selector `json:"selectors"`
}

func (a StructPromptArgs) prompt() string { return "struct_prompt" }

// FixJSONArgs is the arguments for the fix json prompt.
type FixJSONArgs struct {
	JSON string `json:"json"`
}

func (a FixJSONArgs) prompt() string { return "fix_json" }

// PickSelectorArgs is the arguments for the pick selector prompt.
type PickSelectorArgs struct {
	Selectors []master.Selector `json:"selectors"`
	HTML      string            `json:"html"`
	Section   Section           `json:"section"`
}

func (a PickSelectorArgs) prompt() string { return "pick_selector" }
