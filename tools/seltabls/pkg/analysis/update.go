package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"go.lsp.dev/protocol"
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
				s.Documents[string(notification.Params.TextDocument.URI)] = notification.Params.ContentChanges[0].Text
				text := s.Documents[string(notification.Params.TextDocument.URI)]
				comments, err := parsers.ParseStructComments(
					notification.Params.ContentChanges[0].Text,
				)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to get urls and ignores: %w",
						err,
					)
				}
				s.URLs[string(notification.Params.TextDocument.URI)] = append(
					s.URLs[string(notification.Params.TextDocument.URI)],
					comments.URLs...,
				)
				diags, err := GetDiagnosticsForFile(
					ctx,
					&text,
					comments,
				)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to get diagnostics: %w",
						err,
					)
				}
				if len(diags) == 0 {
					return nil, nil
				}
				return &lsp.PublishDiagnosticsNotification{
					Notification: lsp.Notification{
						RPC:    lsp.RPCVersion,
						Method: "textDocument/publishDiagnostics",
					},
					// Params: lsp.PublishDiagnosticsParams{
					//         Diagnostics: diags,
					//         URI:         notification.Params.TextDocument.URI,
					// },
					Params: protocol.PublishDiagnosticsParams{
						URI:         protocol.DocumentURI(notification.Params.TextDocument.URI),
						Version:     uint32(notification.Params.TextDocument.Version),
						Diagnostics: diags,
					},
				}, nil
			}
		}
	}
}
