package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// UpdateDocument updates the state with the given document
func UpdateDocument(
	ctx context.Context,
	notification *lsp.TextDocumentDidChangeNotification,
	documents *safe.Map[uri.URI, string],
	urls *safe.Map[uri.URI, []string],
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
				documents.Set(notification.Params.TextDocument.URI, notification.Params.ContentChanges[0].Text)
				text, ok := documents.Get(notification.Params.TextDocument.URI)
				if !ok {
					return nil, fmt.Errorf("failed to get text")
				}
				comments, err := parsers.ParseStructComments(
					notification.Params.ContentChanges[0].Text,
				)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to get urls and ignores: %w",
						err,
					)
				}
				urls.Set(notification.Params.TextDocument.URI, comments.URLs)
				diags, err := GetDiagnosticsForFile(
					ctx,
					text,
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
					Params: protocol.PublishDiagnosticsParams{
						URI: protocol.DocumentURI(
							notification.Params.TextDocument.URI,
						),
						Version: uint32(
							notification.Params.TextDocument.Version,
						),
						Diagnostics: diags,
					},
				}, nil
			}
		}
	}
}
