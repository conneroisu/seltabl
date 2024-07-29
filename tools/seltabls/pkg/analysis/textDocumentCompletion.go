package analysis

import (
	"bytes"
	"context"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

const (
	maxCompletionLabelLength = 70
	hiddenText               = "..."
)

// CreateTextDocumentCompletion returns the completions for a given text
// document. It checks if the position is within the struct tag and returns the
// selectors if the position is within the struct tag.
//
// It also checks if the position is within the struct tag value and returns
// the selectors if the position is within the struct tag value.
func CreateTextDocumentCompletion(
	ctx context.Context,
	request lsp.TextDocumentCompletionRequest,
	documents *safe.Map[uri.URI, string],
	urls *safe.Map[uri.URI, []string],
	selectors *safe.Map[uri.URI, []master.Selector],
) (response *lsp.TextDocumentCompletionResponse, err error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		response = &lsp.TextDocumentCompletionResponse{
			Response: lsp.Response{
				RPC: lsp.RPCVersion,
				ID:  request.ID,
			},
			Result: []protocol.CompletionItem{},
		}
		content, ok := documents.Get(request.Params.TextDocument.URI)
		if !ok {
			return nil, nil
		}
		selectors, ok := selectors.Get(request.Params.TextDocument.URI)
		if !ok {
			return nil, nil
		}
		url, ok := urls.Get(request.Params.TextDocument.URI)
		if !ok {
			return nil, nil
		}
		log.Debugf("selectors: %v", selectors)
		check, err := CheckPosition(
			request.Params.Position,
			content,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to check position: %w",
				err,
			)
		}
		switch check {
		case parsers.StateInTag:
			for _, key := range completionKeys {
				item := protocol.CompletionItem{
					Label:         key.Label,
					Detail:        key.Detail,
					Documentation: key.Documentation,
					Kind:          protocol.CompletionItemKindField,
				}
				if len(item.Label) > maxCompletionLabelLength {
					item.InsertText = item.Label
					// get the last character to the max length
					item.Label = hiddenText + item.Label[len(item.Label)-maxCompletionLabelLength-len(hiddenText):]
				}
				response.Result = append(
					response.Result,
					item,
				)
			}
		case parsers.StateInTagValue:
			for _, selector := range *selectors {
				item := protocol.CompletionItem{
					Label: selector.Value,
					Detail: fmt.Sprintf(
						"url: '%s'\n Occurances: '%d'\nContext:\n%s",
						(*url)[0],
						selector.Occurances,
						selector.Context,
					),
					CommitCharacters: []string{":", ">", "#"},
					Documentation:    "seltabls",
					Deprecated:       false,
					Kind:             protocol.CompletionItemKindValue,
					InsertTextFormat: protocol.InsertTextFormatPlainText,
					InsertTextMode:   protocol.InsertTextModeAsIs,
				}
				if len(item.Label) > maxCompletionLabelLength {
					item.InsertText = item.Label
					item.Label = hiddenText + item.Label[len(item.Label)-maxCompletionLabelLength-len(hiddenText):]
				}
				response.Result = append(
					response.Result,
					item,
				)
			}
		case parsers.StateAfterColon:
			for _, selector := range *selectors {
				item := protocol.CompletionItem{
					Deprecated: false,
					Detail: fmt.Sprintf(
						"Occurances: '%d'\nContext: \n%s",
						selector.Occurances,
						selector.Context,
					),
					Documentation:    "seltabls",
					CommitCharacters: []string{},
					InsertTextFormat: protocol.InsertTextFormatPlainText,
					InsertTextMode:   protocol.InsertTextModeAsIs,
					Kind:             protocol.CompletionItemKindValue,
					Label: fmt.Sprintf(
						`"%s"`,
						selector.Value,
					),
				}
				if len(item.Label) > maxCompletionLabelLength {
					item.InsertText = item.Label
					item.Label = hiddenText + item.Label[len(item.Label)-maxCompletionLabelLength-len(hiddenText):]
				}
				response.Result = append(
					response.Result,
					item,
				)
			}
		default:
			return nil, nil
		}
		return response, nil
	}
}

// CheckPosition checks if the position is within the struct tag
func CheckPosition(
	position protocol.Position,
	text *string,
) (res parsers.State, err error) {
	fset := token.NewFileSet()
	position.Line = position.Line + 1
	// Parse the source code from a new buffer
	node, err := parser.ParseFile(
		fset,
		"",
		bytes.NewBufferString(*text),
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
		inPosition := parsers.IsPositionInNode(
			structNodes[i],
			position,
			fset,
		)
		inTag := parsers.IsPositionInTag(
			structNodes[i],
			position,
			fset,
		)
		if inPosition && inTag {
			_, inValue := parsers.PositionInStructTagValue(
				structNodes[i],
				position,
				fset,
				text,
			)
			if inValue {
				return parsers.StateInTagValue, nil
			}
			beforeValue := parsers.PositionBeforeValue(
				position,
				text,
			)
			if beforeValue == ':' {
				return parsers.StateAfterColon, nil
			}
			return parsers.StateInTag, nil
		}
	}
	return parsers.StateInvalid, nil
}
