package analysis

import (
	"bytes"
	"context"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"golang.org/x/sync/errgroup"
)

// CreateTextDocumentCompletion returns the completions for a given text document.
// It checks if the position is within the struct tag and returns the selectors
// if the position is within the struct tag.
//
// It also checks if the position is within the struct tag value and returns the selectors
// if the position is within the struct tag value.
func CreateTextDocumentCompletion(
	ctx context.Context,
	s *State,
	request lsp.CompletionRequest,
) (response *lsp.CompletionResponse, err error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		response = &lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  request.ID,
			},
			Result: []lsp.CompletionItem{},
		}
		eg, _ := errgroup.WithContext(ctx)
		eg.Go(func() error {
			// Get the content for the given document.
			content := s.Documents[request.Params.TextDocument.URI]
			// Get the selectors for the given document in current state.
			selectors := s.Selectors[request.Params.TextDocument.URI]
			// Check if the position is within a golang struct tag.
			check, err := s.CheckPosition(request.Params.Position, content)
			if err != nil {
				return nil
			}
			switch check {
			case parsers.StateInTag:
				for _, key := range completionKeys {
					response.Result = append(
						response.Result,
						lsp.CompletionItem{
							Label:         key.Label,
							Detail:        key.Detail,
							Documentation: key.Documentation,
							Kind:          lsp.CompletionKindEnum,
						},
					)
				}
			case parsers.StateInTagValue:
				for _, selector := range selectors {
					response.Result = append(
						response.Result,
						lsp.CompletionItem{
							Label: selector.Value,
							Detail: fmt.Sprintf(
								"Occurances: '%d' \nContext: \n```html\n%s```",
								selector.Occurances,
								selector.Context,
							),
							Documentation: "seltabls",
							Kind:          lsp.CompletionKindReference,
						},
					)
				}
			case parsers.StateAfterColon:
				for _, selector := range selectors {
					response.Result = append(
						response.Result,
						lsp.CompletionItem{
							Label: "\"" + selector.Value + "\"",
							Detail: fmt.Sprintf(
								"Occurances: '%d' \nContext: \n```html\n%s```",
								selector.Occurances,
								selector.Context,
							),
							Documentation: "seltabls",
							Kind:          lsp.CompletionKindReference,
						},
					)
				}
			default:
				return nil
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			return response, fmt.Errorf("failed to get completions: %w", err)
		}
		return response, nil
	}
}

// CheckPosition checks if the position is within the struct tag
func (s *State) CheckPosition(
	position lsp.Position,
	text string,
) (res parsers.State, err error) {
	var inValue bool
	// Create a new token file set
	fset := token.NewFileSet()
	position.Line = position.Line + 1
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
	for i := range structNodes {
		// Check if the position is within the struct node
		inPosition := parsers.IsPositionInNode(structNodes[i], position, fset)
		// Check if the position is within a struct tag
		inTag := parsers.IsPositionInTag(structNodes[i], position, fset)
		if inPosition && inTag {
			// Check if the position is within a struct tag value (i.e. value inside and including " and " characters)
			_, inValue = parsers.PositionInStructTagValue(
				structNodes[i],
				position,
				fset,
			)
			if inValue {
				return parsers.StateInTagValue, nil
			}
			if parsers.PositionBeforeValue(position, text) == '"' {
				// If the position is before a double quote, return the state in the tag Value
				// Also return the key of the struct tag before the double quote aka our position.
				// TODO: Get the key of the struct tag before the double quote
				return parsers.StateInTagValue, nil
			}
			// If we are in the tag, we should return completion items for the struct tag
			// that are not yet set/defined
			return parsers.StateInTag, nil
		}
	}
	return parsers.StateInvalid, nil
}
