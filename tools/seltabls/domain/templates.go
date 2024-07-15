package domain

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

//go:embed templates.tmpl
var templatesTmpl string

// Template is the must compile template
var Template = template.Must(template.New("templates").Parse(templatesTmpl))

// prompter is an interface for prompting it provided the template name to use.
type prompter interface {
	prompt() string
}

// NewPrompt creates a new prompt for the given args
func NewPrompt(
	promptArgs prompter,
) (string, error) {
	name := promptArgs.prompt()
	buf := new(bytes.Buffer)
	err := Template.ExecuteTemplate(buf, name, promptArgs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type sectionErrorArgs struct {
	Error error
}

func (a sectionErrorArgs) prompt() string {
	return "section_error"
}

// SectionAggregateArgs is the arguments for the section aggregate prompt.
type SectionAggregateArgs struct {
	Structs []string
	Content string
}

func (a SectionAggregateArgs) prompt() string {
	return "section_aggregate"
}

// IdentifyErrorArgs is the arguments for the identify error prompt.
type IdentifyErrorArgs struct {
	Error error
}

func (a IdentifyErrorArgs) prompt() string { return "identify_error" }

// StructAggregateArgs is the arguments for the struct aggregate prompt.
type StructAggregateArgs struct {
	Selectors []master.Selector
	Content   string
	Schemas   []string
}

func (a StructAggregateArgs) prompt() string { return "struct_aggregate" }

// IdentifyPromptArgs is the arguments for the identify prompt.
type IdentifyPromptArgs struct {
	URL     string
	Content string
}

func (a IdentifyPromptArgs) prompt() string { return "identify_prompt" }

// IdentifyAggregateArgs is the arguments for the identify aggregate prompt.
type IdentifyAggregateArgs struct {
	Content string
	Schemas []string
}

func (a IdentifyAggregateArgs) prompt() string { return "identify_aggregate" }

// StructPromptArgs is the arguments for the struct prompt.
type StructPromptArgs struct {
	URL       string
	Content   string
	Selectors []master.Selector
}

func (a StructPromptArgs) prompt() string { return "struct_prompt" }

// TestFilePromptArgs is the arguments for the test file prompt.
type TestFilePromptArgs struct {
	Version     string
	Name        string
	URL         string
	PackageName string
}

func (a TestFilePromptArgs) prompt() string { return "test_file" }

// StructFilePromptArgs is the arguments for the struct file prompt.
type StructFilePromptArgs struct {
	Name           string
	URL            string
	IgnoreElements []string
	Fields         []Field
	PackageName    string
}

func (a StructFilePromptArgs) prompt() string { return "struct_file" }
