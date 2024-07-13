package domain

import (
	"bytes"
	_ "embed"
	"fmt"
	"reflect"
	"text/template"
)

var (
	//go:embed prompts.tmpl
	promptsTemplate string
	// PromptsTemplate is the template for the prompts file.
	PromptsTemplate *template.Template = template.Must(template.New("prompts_template").Parse(promptsTemplate))
)

// NewPrompt creates a new filled out template for a prompt from the
// `struct.tmpl` template.
//
// It takes as input a interface that implements the prompter interface.
// That is, it has a method called Prompt() that returns a string.
func NewPrompt[
	T NewAggregateStuctPromptArgs | NewErrorPromptArgs | NewStructFileArgs | NewErrorAggregatePromptArgs | NewStructContentArgs | NewSelectorPromptArgs,
](
	args T,
) (buf bytes.Buffer, err error) {
	t := reflect.TypeOf(args)
	field := t.FieldByIndex([]int{0})
	key := field.Tag.Get("prompt")
	fmt.Println("Tag value (Method 1):", key)
	err = PromptsTemplate.ExecuteTemplate(&buf, key, args)
	if err != nil {
		return buf, fmt.Errorf(
			"failed to execute struct file template: %w",
			err,
		)
	}
	return buf, nil
}
