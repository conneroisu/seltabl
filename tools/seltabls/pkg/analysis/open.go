package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"golang.org/x/sync/errgroup"
)

// OpenDocument opens a document in the state and returns any diagnostics for the document
//
// uri is the uri of the document
//
// content is the content of the document
func (s *State) OpenDocument(
	ctx context.Context,
	req lsp.NotificationDidOpenTextDocument,
) (response *lsp.PublishDiagnosticsNotification, err error) {
	diags := []lsp.Diagnostic{}
	response = &lsp.PublishDiagnosticsNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/publishDiagnostics",
		},
		Params: lsp.PublishDiagnosticsParams{
			URI:         req.Params.TextDocument.URI,
			Diagnostics: diags,
		},
	}
	uri := req.Params.TextDocument.URI
	eg, ctx := errgroup.WithContext(ctx)
	s.Documents[uri] = req.Params.TextDocument.Text
	data, err := parsers.ParseStructComments(
		req.Params.TextDocument.Text,
	)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
		return response, nil
	}
	for _, url := range data.URLs {
		eg.Go(func() error {
			s.URLs[uri] = append(s.URLs[uri], url)
			s.Selectors[uri], err = GetSelectors(
				ctx,
				s,
				url,
				data.IgnoreElements,
			)
			if err != nil {
				return fmt.Errorf(
					"failed to get selectors for url (%s): %w",
					url,
					err,
				)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get selectors for urls: %w", err)
	}
	response = &lsp.PublishDiagnosticsNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/publishDiagnostics",
		},
		Params: lsp.PublishDiagnosticsParams{
			URI:         req.Params.TextDocument.URI,
			Diagnostics: diags,
		},
	}
	diags, err = s.GetDiagnosticsForFile(
		&req.Params.TextDocument.Text,
		data,
	)
	if err != nil {
		s.Logger.Printf("failed to get diagnostics for file: %s\n", err)
	}
	return response, nil
}
