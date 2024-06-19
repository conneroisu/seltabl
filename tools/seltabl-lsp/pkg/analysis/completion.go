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
	// headerTag is the tag used to match a header cell's Value.
	headerTag = lsp.CompletionItem{
		Label:         "seltabl",
		Detail:        "Title Text for the header",
		Documentation: "This is the documentation for the header",
	}
	// selectorDataTag is the tag used to mark a data cell.
	selectorDataTag = lsp.CompletionItem{
		Label:         "dSel",
		Detail:        "Title Text for the data selector",
		Documentation: "This is the documentation for the data selector",
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorHeaderTag = lsp.CompletionItem{
		Label:         "hSel",
		Detail:        "Title Text for the header selector",
		Documentation: "This is the documentation for the header selector",
	}
	// selectorQueryTag is the tag used to signify selecting aspects of a cell
	selectorQueryTag = lsp.CompletionItem{
		Label:         "qSel",
		Detail:        "Title Text for the query selector",
		Documentation: "This is the documentation for the query selector",
	}
	// selectorMustBePresentTag is the tag used to signify selecting aspects of a cell
	selectorMustBePresentTag = lsp.CompletionItem{
		Label:         "must",
		Detail:        "Title Text for the must be present selector",
		Documentation: "This is the documentation for the must be present selector",
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorControlTag = lsp.CompletionItem{
		Label:         "ctl",
		Detail:        "Title Text for the control selector",
		Documentation: "This is the documentation for the control selector",
	}
	// keys is the slice of keys to return for completions inside a struct tag but not a "" selector
	keys = []lsp.CompletionItem{
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
	isPositionInStructTag, err := s.CheckPosition(*pos, text)
	if err != nil {
		return nil, fmt.Errorf("failed to check position: %w", err)
	}
	if isPositionInStructTag != nil {
		s.Logger.Println("isPositionInStructTag", *isPositionInStructTag)
		if *isPositionInStructTag {
			for _, selector := range selectors {
				items = append(items, lsp.CompletionItem{
					Label:         selector.Value,
					Detail:        "context: \n" + selector.Context,
					Documentation: "seltabl-lsp",
				})
			}
		} else {

			for _, key := range keys {
				items = append(items, lsp.CompletionItem{
					Label:         key.Label,
					Detail:        key.Detail,
					Documentation: key.Documentation,
				})
			}

		}
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
func (s *State) CheckPosition(position lsp.Position, text string) (res *bool, err error) {
	var TRUE = true
	var FALSE = false
	// Read the Go source code from a file
	sourceCode := bytes.NewBufferString(text)
	// Create a new token file set
	fset := token.NewFileSet()
	// Parse the source code
	node, err := parser.ParseFile(fset, "", sourceCode, parser.Trace)
	if err != nil {
		return nil, err
	}
	// Find the struct node in the AST
	structNodes := parsers.FindStructNodes(node)
	for _, structNode := range structNodes {
		// Check if the position is within the struct node
		if parsers.IsPositionInNode(structNode, position, fset) {
			// Check if the position is within a struct tag
			if parsers.IsPositionInTag(structNode, position, fset) {
				// Check if the position is within a struct tag value (i.e. value inside and including " and " characters)
				if parsers.IsPositionInStructTagValue(structNode, position, fset) {
					return &TRUE, nil
				}
			}
			return &FALSE, nil
		}
	}
	return nil, nil
}
