package domain

import (
	"bytes"
	_ "embed"
	"fmt"
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

// Responder is an interface for responding to a given prompt.
type Responder interface{ respond() string }

func (i IdentifyResponse) respond() string {
	return fmt.Sprintf(
		"Sections:\n%s\nPackageName: %s",
		i.Sections,
		i.PackageName,
	)
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
	Structs   []string          `json:"structs"`   // required
	Content   string            `json:"content"`   // required
	Selectors []master.Selector `json:"selectors"` // required
}

func (a SectionAggregateArgs) prompt() string {
	return "section_aggregate"
}

// IdentifyErrorArgs is the arguments for the identify error prompt.
type IdentifyErrorArgs struct {
	Error error `json:"error"` // required
}

func (a IdentifyErrorArgs) prompt() string { return "section_error" }

// StructAggregateArgs is the arguments for the struct aggregate prompt.
type StructAggregateArgs struct {
	Selectors []master.Selector // required
	Content   string            // required
	Schemas   []string          // required
}

func (a StructAggregateArgs) prompt() string { return "struct_aggregate" }

// IdentifyPromptArgs is the arguments for the identify prompt.
type IdentifyPromptArgs struct {
	URL         string `json:"url"`
	Content     string `json:"content"`
	NumSections int    `json:"num-sections"`
}

func (a IdentifyPromptArgs) prompt() string { return "identify_prompt" }

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

// TestFilePromptArgs is the arguments for the test file prompt.
type TestFilePromptArgs struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	PackageName string `json:"package-name"`
}

func (a TestFilePromptArgs) prompt() string { return "test_file" }

// StructFilePromptArgs is the arguments for the struct file prompt.
type StructFilePromptArgs struct {
	Name           string   `json:"name,omitempty"`
	URL            string   `json:"url,omitempty"`
	IgnoreElements []string `json:"ignore-elements,omitempty"`
	Fields         []Field  `json:"fields"`
	PackageName    string   `json:"package-name"`
}

func (a StructFilePromptArgs) prompt() string { return "struct_file" }

// DecodeErrorArgs is the arguments for the decode error prompt.
type DecodeErrorArgs struct {
	Error   error  `json:"error"`
	Message string `json:"message"`
}

func (a DecodeErrorArgs) prompt() string { return "decode_error" }
