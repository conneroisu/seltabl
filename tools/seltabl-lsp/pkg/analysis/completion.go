package analysis

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/parsers"
)

var (
	// completionKeys is the slice of completionKeys to return for completions inside a struct tag but not a "" selector
	completionKeys = []lsp.CompletionItem{
		headerTag,
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
	document *lsp.TextDocumentIdentifier,
	pos *lsp.Position,
) (response *lsp.CompletionResponse, err error) {
	text := s.Documents[document.URI]
	selectors := s.Selectors[document.URI]
	items := []lsp.CompletionItem{}
	check, err := s.CheckPosition(*pos, text)
	if err != nil {
		return nil, fmt.Errorf("failed to check position: %w", err)
	}
	switch check {
	case parsers.StateInTag:
		s.Logger.Println("Found position in struct tag")
		for _, key := range completionKeys {
			items = append(items, lsp.CompletionItem{
				Label:         key.Label,
				Detail:        key.Detail,
				Documentation: key.Documentation,
				Kind:          lsp.Enum,
			})
		}
	case parsers.StateInTagValue:
		s.Logger.Println("Found position in struct tag value")
		for _, selector := range selectors {
			items = append(items, lsp.CompletionItem{
				Label:         selector.Value,
				Detail:        "context: \n" + selector.Context,
				Documentation: "seltabl-lsp",
				Kind:          lsp.Reference,
			})
		}
	case parsers.StateAfterColon:
		s.Logger.Println("Found position in struct tag after colon")
		for _, selector := range selectors {
			items = append(items, lsp.CompletionItem{
				Label:         "\"" + selector.Value + "\"",
				Detail:        "context: \n" + selector.Context,
				Documentation: "seltabl-lsp",
				Kind:          lsp.Reference,
			})
		}
	case parsers.StateInvalid:
	default:
		return nil, nil
	}
	return &lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: items,
	}, nil
}

// CheckPosition checks if the position is within the struct tag
func (s *State) CheckPosition(
	position lsp.Position,
	text string,
) (res parsers.State, err error) {
	var inValue bool
	// Read the Go source code from a file
	sourceCode := bytes.NewBufferString(text)
	// Create a new token file set
	fset := token.NewFileSet()
	// Parse the source code
	node, err := parser.ParseFile(fset, "", sourceCode, parser.Trace)
	if err != nil {
		return parsers.StateInvalid, fmt.Errorf(
			"failed to parse struct: %w",
			err,
		)
	}
	// Find the struct node in the AST
	structNodes := parsers.FindStructNodes(node)
	for _, structNode := range structNodes {
		// Check if the position is within the struct node
		inPosition := parsers.IsPositionInNode(structNode, position, fset)
		// Check if the position is within a struct tag
		inTag := parsers.IsPositionInTag(structNode, position, fset)
		if inPosition && inTag {
			// Check if the position is within a struct tag value (i.e. value inside and including " and " characters)
			_, inValue = parsers.IsPositionInStructTagValue(
				structNode,
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
