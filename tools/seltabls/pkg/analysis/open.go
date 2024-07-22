package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"go.lsp.dev/protocol"
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
	s *State,
	req lsp.NotificationDidOpenTextDocument,
) (*lsp.PublishDiagnosticsNotification, error) {
	select {
	case <-(ctx).Done():
		return nil, fmt.Errorf("context cancelled: %w", (ctx).Err())
	default:
		eg, ctx := errgroup.WithContext(ctx)
		uri := req.Params.TextDocument.URI
		s.Documents[string(uri)] = req.Params.TextDocument.Text
		data, err := parsers.ParseStructComments(req.Params.TextDocument.Text)
		if err != nil {
			return nil, nil
		}
		s.URLs[string(uri)] = append(s.URLs[string(uri)], data.URLs...)
		for _, url := range data.URLs {
			eg.Go(func() error {
				s.Selectors[string(uri)], err = parsers.GetSelectors(
					ctx,
					&s.Database,
					url,
					data.IgnoreElements,
					2,
				)
				if err != nil {
					return err
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
			&req.Params.TextDocument.Text,
			data,
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
				URI:         protocol.DocumentURI(req.Params.TextDocument.URI),
				Diagnostics: diags,
			},
		}, nil
	}
}
