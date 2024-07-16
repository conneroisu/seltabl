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
	notification *lsp.TextDocumentDidChangeNotification,
) (*lsp.PublishDiagnosticsNotification, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
			default:
				s.Documents[notification.Params.TextDocument.URI] = notification.Params.ContentChanges[len(notification.Params.ContentChanges)-1].Text
				data, err := parsers.ParseStructComments(
					notification.Params.ContentChanges[len(notification.Params.ContentChanges)-1].Text,
				)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to get urls and ignores: %w",
						err,
					)
				}
				s.URLs[notification.Params.TextDocument.URI] = append(
					s.URLs[notification.Params.TextDocument.URI],
					data.URLs...,
				)
				var ds []lsp.Diagnostic
				for i := range notification.Params.ContentChanges {
					diags, err := GetDiagnosticsForFile(
						ctx,
						s,
						&notification.Params.ContentChanges[i].Text,
						data,
					)
					if err != nil {
						return nil, fmt.Errorf(
							"failed to get diagnostics for file: %w",
							err,
						)
					}
					ds = append(
						ds,
						diags...)
				}
				return &lsp.PublishDiagnosticsNotification{
					Notification: lsp.Notification{
						RPC:    "2.0",
						Method: "textDocument/publishDiagnostics",
					},
					Params: lsp.PublishDiagnosticsParams{
						Diagnostics: []lsp.Diagnostic{},
						URI:         notification.Params.TextDocument.URI,
					},
				}, nil
			}
		}
	}
}
