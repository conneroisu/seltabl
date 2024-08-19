package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
	"golang.org/x/sync/errgroup"
)

// OpenDocument opens a document in the state and returns any diagnostics for the document.
//
// uri is the uri of the document.
//
// content is the content of the document.
//
// On the opening of any document, the state is updated with the content of the document
// and the diagnostics for the document are returned.
func OpenDocument(
	ctx context.Context,
	req *lsp.NotificationDidOpenTextDocument,
	db *data.Database[master.Queries],
	documents *safe.Map[uri.URI, *parsers.GoFile],
) (*lsp.PublishDiagnosticsNotification, error) {
	for {
		select {
		case <-(ctx).Done():
			return nil, ctx.Err()
		default:
			eg, ctx := errgroup.WithContext(ctx)
			source, err := parsers.ParseSource(
				req.Params.TextDocument.Text,
				req.Params.TextDocument.URI.Filename(),
				true,
			)
			documents.Set(
				req.Params.TextDocument.URI,
				source,
			)
			for _, structure := range source.Structs {
				eg.Go(func() error {
					var occur int
					if len(structure.Fields) <= 0 {
						return nil
					}
					var sels []master.Selector
					for _, ele := range structure.SeltablIgnores {
						sels, err = parsers.GetSelectors(
							ctx,
							db,
							url,
							ele,
							occur,
						)
						if err != nil {
							return err
						}
						selectors.Set(
							req.Params.TextDocument.URI,
							append(sels, sels...),
						)
					}
					return nil
				})
			}
			err = eg.Wait()
			if err != nil {
				return nil, fmt.Errorf(
					"failed to get selectors for urls: %w",
					err,
				)
			}
			diags, err := GetDiagnosticsForFile(
				ctx,
				db,
				(*string)(&req.Params.TextDocument.URI),
				documents,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to get diagnostics for file: %w",
					err,
				)
			}
			if len(diags) == 0 {
				return nil, nil
			}
			return &lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    lsp.RPCVersion,
					Method: string(methods.NotificationPublishDiagnostics),
				},
				Params: protocol.PublishDiagnosticsParams{
					URI: protocol.DocumentURI(
						req.Params.TextDocument.URI,
					),
					Diagnostics: diags,
				},
			}, nil
		}
	}
}
