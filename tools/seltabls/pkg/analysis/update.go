package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// UpdateDocument updates the state with the given document.
//
// Upon updating the document, it returns new diagnostics for the document.
func UpdateDocument(
	ctx context.Context,
	notification *lsp.TextDocumentDidChangeNotification,
	db *data.Database[master.Queries],
	documents *safe.Map[uri.URI, string],
	urls *safe.Map[uri.URI, []string],
	selectors *safe.Map[uri.URI, []master.Selector],
) (*lsp.PublishDiagnosticsNotification, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			documents.Set(
				notification.Params.TextDocument.URI,
				notification.Params.ContentChanges[0].Text,
			)
			comments, err := parsers.ParseStructComments(
				notification.Params.ContentChanges[0].Text,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to get urls and ignores: %w",
					err,
				)
			}
			urls.Set(
				notification.Params.TextDocument.URI,
				comments.URLs,
			)
			for i, url := range comments.URLs {
				if len(comments.URLs) == 0 {
					return nil, fmt.Errorf("no urls found")
				}
				if len(comments.IgnoreElements) == 0 {
					return nil, fmt.Errorf("no ignore elements found")
				}
				if len(comments.Occurrences) == 0 {
					return nil, fmt.Errorf("no occurrences found")
				}
				sels, err := parsers.GetSelectors(
					ctx,
					db,
					url,
					comments.IgnoreElements[i],
					comments.Occurrences[i],
				)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to get selectors: %w",
						err,
					)
				}
				selectors.Set(notification.Params.TextDocument.URI, sels)
			}
			diags, err := GetDiagnosticsForFile(
				ctx,
				&notification.Params.ContentChanges[0].Text,
				&comments,
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
