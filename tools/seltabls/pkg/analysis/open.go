package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
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
) (response rpc.MethodActor, err error) {
	response = &lsp.Response{
		RPC: lsp.RPCVersion,
		ID:  *req.ID,
	}
	select {
	case <-(ctx).Done():
		return response, fmt.Errorf("context cancelled: %w", (ctx).Err())
	default:
		eg, ctx := errgroup.WithContext(ctx)
		uri := req.Params.TextDocument.URI
		s.Documents[uri] = req.Params.TextDocument.Text
		data, err := parsers.ParseStructComments(req.Params.TextDocument.Text)
		if err != nil {
			return response, fmt.Errorf("failed to get selectors for urls: %w", err)
		}
		s.URLs[uri] = append(s.URLs[uri], data.URLs...)
		for _, url := range data.URLs {
			eg.Go(func() error {
				s.Selectors[uri], err = parsers.GetSelectors(
					ctx,
					&s.Database,
					url,
					data.IgnoreElements,
				)
				if err != nil {
					return err
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return response, fmt.Errorf(
				"failed to get selectors for urls: %w",
				err,
			)
		}
		var diags []lsp.Diagnostic
		diags, err = GetDiagnosticsForFile(
			ctx,
			s,
			&req.Params.TextDocument.Text,
			data,
		)
		if err != nil {
			return response, fmt.Errorf("failed to get diagnostics for file: %w", err)
		}
		response = lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    lsp.RPCVersion,
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         req.Params.TextDocument.URI,
				Diagnostics: diags,
			},
		}
		return response, nil
	}
}
