package domain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

// TestNewAggregatePrompt tests the NewAggregatePrompt struct
func TestNewAggregatePrompt(t *testing.T) {
	a := assert.New(t)
	content, err := NewPrompt(
		PromptAggregateSections{
			Structs: []string{"ex json 1 ", "ex json 2 ", "ex json 3 "},
			Content: `<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>`,
			Selectors: []master.Selector{
				{
					ID:         2,
					Value:      "dsaf",
					UrlID:      2,
					Occurances: 2,
					Context:    "<html>",
				},
			},
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}

func TestIdentifyAggregateArgs(t *testing.T) {
	a := assert.New(t)
	content, err := NewPrompt(
		IdentifyAggregateArgs{
			Content: "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>",
			Schemas: []string{"dsafsd", "dsazfdasdfasf"},
			Selectors: []master.Selector{
				{
					ID:         2,
					Value:      "dsaf",
					UrlID:      2,
					Occurances: 2,
					Context:    "<html>",
				},
			},
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}

// TestStructPromptArgs tests the StructPromptArgs struct.
func TestStructPromptArgs(t *testing.T) {
	a := assert.New(t)
	content, err := NewPrompt(
		StructPromptArgs{
			URL:     "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
			Content: "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>",
			Selectors: []master.Selector{
				{
					ID:         1,
					Value:      "html > body > table > tbody > tr > td:nth-child(1)",
					Occurances: 1,
					Context:    "<html></html>",
				},
			},
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}

// TestSectionErrorPromptArgs tests the SectionErrorPromptArgs struct.
func TestSectionErrorPromptArgs(t *testing.T) {
	a := assert.New(t)
	content, err := NewPrompt(
		PromptBetterError{
			Error: errors.New(
				"failed to get the content of the url: https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
			),
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}

// TestNewPromptIdentifyArgs tests the IdentifyPromptArgs struct with a single selector.
func TestNewPromptIdentifyArgs(t *testing.T) {
	a := assert.New(t)
	content, err := NewPrompt(
		IdentifyArgs{
			URL:         "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
			Content:     "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>",
			NumSections: 3,
			Selectors: []master.Selector{
				{
					ID:         1,
					Value:      "html > body > table > tbody > tr > td:nth-child(1)",
					Occurances: 1,
					Context:    "<html></html>",
				},
			},
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}

// TestNewPromptPickSelector tests the NewPrompt function with a PickSelectorArgs struct.
func TestNewPromptPickSelector(t *testing.T) {
	a := assert.New(t)
	content, err := NewPrompt(
		PickSelectorArgs{
			Selectors: []master.Selector{
				{
					Value: "html > body > table#dataTable > tr:nth-child(1) > td:nth-child(1)",
				},
			},
			HTML: "<html><body><table id=\"dataTable\"><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>",
			Section: Section{
				Name:        "Test",
				Description: "Test Section",
				CSS:         "html > body > table#dataTable > tr:nth-child(1) > td:nth-child(1)",
			},
		},
	)

	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}

// TestNewPromptBetterError tests the NewStructFileContent struct
func TestNewPromptBetterError(t *testing.T) {
	a := assert.New(t)
	content, err := NewPrompt(
		PromptBetterError{
			Error: fmt.Errorf(
				"failed to parse struct: failed to get data rows html: failed to get html: failed to get doc: open /Users/hsz/Projects/github.com/conneroisu/seltabl/testdata/ab_num_table.html: no such file or directory",
			),
			History: []anthropic.Message{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: []anthropic.MessageContent{anthropic.NewTextMessageContent("foo")},
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: []anthropic.MessageContent{anthropic.NewTextMessageContent("bar")},
				},
			},
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}
