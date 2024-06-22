package analysis

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
)

var (
	// completionKeys is the slice of completionKeys to return for completions inside a struct tag but not a "" selector
	completionKeys = []lsp.CompletionItem{
		selectorDataTag,
		selectorHeaderTag,
		selectorQueryTag,
		selectorMustBePresentTag,
		selectorControlTag,
	}
)

// CreateTextDocumentCompletion returns the completions for a given text document.
// It checks if the position is within the struct tag and returns the selectors
// if the position is within the struct tag.
//
// It also checks if the position is within the struct tag value and returns the selectors
// if the position is within the struct tag value.
func (s *State) CreateTextDocumentCompletion(
	id int,
	document lsp.TextDocumentIdentifier,
	pos lsp.Position,
) (response lsp.CompletionResponse, err error) {
	response.Response = lsp.Response{
		RPC: "2.0",
		ID:  id,
	}
	response.Result = []lsp.CompletionItem{}
	// Get the content for the given document.
	content := s.Documents[document.URI]
	// Get the selectors for the given document in current state.
	selectors := s.Selectors[document.URI]
	// Check if the position is within a golang struct tag.
	check, err := s.CheckPosition(pos, content)
	if err != nil {
		return lsp.CompletionResponse{}, fmt.Errorf(
			"failed to check position: %w",
			err,
		)
	}
	switch check {
	case parsers.StateInTag:
		for _, key := range completionKeys {
			response.Result = append(response.Result, lsp.CompletionItem{
				Label:         key.Label,
				Detail:        key.Detail,
				Documentation: key.Documentation,
				Kind:          lsp.Enum,
			})
		}
	case parsers.StateInTagValue:
		for _, selector := range selectors {
			response.Result = append(response.Result, lsp.CompletionItem{
				Label:         selector.Value,
				Detail:        "context: \n" + selector.Context,
				Documentation: "seltabls",
				Kind:          lsp.Reference,
			})
		}
	case parsers.StateAfterColon:
		for _, selector := range selectors {
			response.Result = append(response.Result, lsp.CompletionItem{
				Label:         "\"" + selector.Value + "\"",
				Detail:        "context: \n" + selector.Context,
				Documentation: "seltabls",
				Kind:          lsp.Reference,
			})
		}
	case parsers.StateInvalid:
		return response, nil
	default:
		return response, nil
	}
	return response, nil
}

// CheckPosition checks if the position is within the struct tag
func (s *State) CheckPosition(
	position lsp.Position,
	text string,
) (res parsers.State, err error) {
	var inValue bool
	// Create a new token file set
	fset := token.NewFileSet()
	// Parse the source code from a new buffer
	node, err := parser.ParseFile(
		fset,
		"",
		bytes.NewBufferString(text),
		parser.Trace,
	)
	if err != nil {
		return parsers.StateInvalid,
			fmt.Errorf("failed to parse struct: %w", err)
	}
	// Find the struct node in the AST
	structNodes := parsers.FindStructNodes(node)
	// for each of the nodes
	for i := range structNodes {
		// Check if the position is within the struct node
		inPosition := parsers.IsPositionInNode(structNodes[i], position, fset)
		// Check if the position is within a struct tag
		inTag := parsers.IsPositionInTag(structNodes[i], position, fset)
		if inPosition && inTag {
			// Check if the position is within a struct tag value (i.e. value inside and including " and " characters)
			_, inValue = parsers.IsPositionInStructTagValue(
				structNodes[i],
				position,
				fset,
			)
			if inValue {
				return parsers.StateInTagValue, nil
			}
			// Check if the position is at / after a colon
			if parsers.PositionBeforeValue(position, text) == ':' {
				return parsers.StateAfterColon, nil
			}
			if parsers.PositionBeforeValue(position, text) == '"' {
				return parsers.StateInTagValue, nil
			}
			return parsers.StateInTag, nil
		}
	}
	return parsers.StateInvalid, nil
}
