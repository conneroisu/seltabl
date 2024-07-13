package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type args[
	T NewAggregateStuctPromptArgs | NewErrorPromptArgs | NewStructFileArgs | NewErrorAggregatePromptArgs | NewStructContentArgs | NewSelectorPromptArgs,
] struct {
	args T
}

var (
	// newStructContentArgs is a slice of test cases
	newStructContentArgs = []args[NewStructContentArgs]{{args: NewStructContentArgs{URL: "https://github.com", Name: "GithubPage", IgnoreElements: []string{"script", "style", "link", "img", "footer", "header"}, Fields: []Field{{Name: "A", Type: "string", Description: "A description of the field", HeaderSelector: "tr:nth-child(1) td:nth-child(1)", DataSelector: "tr td:nth-child(1)", ControlSelector: "$text", MustBePresent: "NCAA Codes"}}}}}
	// newAggregateStuctPromptArgs is a slice of test cases
	newAggregateStuctPromptArgs = []args[NewAggregateStuctPromptArgs]{{args: NewAggregateStuctPromptArgs{URL: "https://github.com", HTMLContent: "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>", Selectors: []string{"html", "html > body", "html > body > table", "html > body > table > tbody", "html > body > table > tbody > tr", "html > body > table > tbody > tr > td", "html > body > table > tbody > tr > td > a[href]"}, Schemas: []string{`{ "fields": [ { "name": "A", "type": "string", "description": "A description of the field", "header-selector": "tr:nth-child(1) td:nth-child(1)", "data-selector": "tr td:nth-child(1)", "control-selector": "$text", "must-be-present": "NCAA Codes" }, { "name": "B", "type": "int		", "description": "A description of the field", "header-selector": "tr:nth-child(1) td:nth-child(2)", "data-selector": "tr td:nth-child(2)", "control-selector": "$text", "must-be-present": "NCAA Codes" } ] }`}}}}
	// newSelectorPromptArgs is a slice of test cases
	newSelectorPromptArgs = []args[NewSelectorPromptArgs]{{args: NewSelectorPromptArgs{SelectorName: "header-selector", SelectorDescription: "used to find the header row and column for the field in the given struct.", URL: "https://github.com", Content: "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>", Selectors: []string{"html", "html > body", "html > body > table", "html > body > table > tbody", "html > body > table > tbody > tr", "html > body > table > tbody > tr > td", "html > body > table > tbody > tr > td > a[href]"}}}}
	// newErrorPromptArgs is a slice of test cases
	newErrorPromptArgs = []args[NewErrorPromptArgs]{{args: NewErrorPromptArgs{Error: fmt.Errorf("failed to get the content of the url: Get https://github.com: dial tcp: lookup github.com on 127.0.0.53:53: read udp 127.0.0.1:53448->127.0.0.53:53: i/o timeout"), Out: "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>"}}}
	// newErrorAggregatePromptArgs is a slice of test cases
	newErrorAggregatePromptArgs = []args[NewErrorAggregatePromptArgs]{{args: NewErrorAggregatePromptArgs{Errors: []error{fmt.Errorf("failed to get the content of the url: Get https://github.com: dial tcp: lookup github.com on 127.0.0.53:53: read udp 127.0.0.1:53448->127.0.0.53:53: i/o timeout")}, Out: "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>"}}}
	// newStructFileArgs is a slice of test cases
	newStructFileArgs = []args[NewStructFileArgs]{{args: NewStructFileArgs{URL: "https://github.com", PackageName: "github"}}}
)

// TestNewPrompt tests the NewPrompt function.
func TestNewPrompt(t *testing.T) {
	t.Run("Test NewAggregateStuctPrompt", func(t *testing.T) {
		for _, arg := range newAggregateStuctPromptArgs {
			arg.RunTestCase(t)
		}
	})
	t.Run("Test NewErrorPrompt", func(t *testing.T) {
		for _, arg := range newErrorPromptArgs {
			arg.RunTestCase(t)
		}
	})
	t.Run("Test NewErrorAggregatePrompt", func(t *testing.T) {
		for _, arg := range newErrorAggregatePromptArgs {
			arg.RunTestCase(t)
		}
	})
	t.Run("Test NewSelectorPrompt", func(t *testing.T) {
		for _, arg := range newSelectorPromptArgs {
			arg.RunTestCase(t)
		}
	})
	t.Run("Test NewStructContent", func(t *testing.T) {
		for _, arg := range newStructContentArgs {
			arg.RunTestCase(t)
		}
	})
	t.Run("Test NewStructFile", func(t *testing.T) {
		for _, arg := range newStructFileArgs {
			arg.RunTestCase(t)
		}
	})
}

// RunTestCase runs a test case for the given args.
func (a *args[T]) RunTestCase(t *testing.T) {
	got, err := NewPrompt(a.args)
	assert.NoError(t, err)
	t.Logf("got: \nvalue: \n%s", got.String())
}
