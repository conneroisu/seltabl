package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
)

// UpdateDocument updates the state with the given document
func UpdateDocument(
	ctx context.Context,
	s *State,
	request *lsp.TextDocumentDidChangeNotification,
) (response *lsp.PublishDiagnosticsNotification, err error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			response = &lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					Diagnostics: []lsp.Diagnostic{},
					URI:         request.Params.TextDocument.URI,
				},
			}
			select {
			case <-ctx.Done():
				return response, fmt.Errorf("context cancelled: %w", ctx.Err())
			default:
				s.Documents[request.Params.TextDocument.URI] = request.Params.ContentChanges[0].Text
				data, err := parsers.ParseStructComments(
					request.Params.ContentChanges[0].Text,
				)
				if err != nil {
					return response, fmt.Errorf(
						"failed to get urls and ignores: %w",
						err,
					)
				}
				s.URLs[request.Params.TextDocument.URI] = append(
					s.URLs[request.Params.TextDocument.URI],
					data.URLs...,
				)
				for i := range request.Params.ContentChanges {
					diags, err := GetDiagnosticsForFile(
						ctx,
						s,
						&request.Params.ContentChanges[i].Text,
						data,
					)
					if err != nil {
						return response, fmt.Errorf(
							"failed to get diagnostics for file: %w",
							err,
						)
					}
					response.Params.Diagnostics = append(
						response.Params.Diagnostics,
						diags...)
				}
				return response, nil
			}
		}
	}
}
