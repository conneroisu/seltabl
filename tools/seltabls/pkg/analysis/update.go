package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"golang.org/x/sync/errgroup"
)

// UpdateDocument updates the state with the given document
func (s *State) UpdateDocument(
	ctx context.Context,
	request *lsp.TextDocumentDidChangeNotification,
) (response lsp.PublishDiagnosticsNotification, err error) {
	var eg *errgroup.Group
	response = lsp.PublishDiagnosticsNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/publishDiagnostics",
		},
		Params: lsp.PublishDiagnosticsParams{
			URI:         request.Params.TextDocument.URI,
			Diagnostics: []lsp.Diagnostic{},
		},
	}
	eg, _ = errgroup.WithContext(ctx)
	eg.Go(func() error {
		s.Documents[request.Params.TextDocument.URI] = request.Params.ContentChanges[0].Text
		data, err := parsers.ParseStructComments(request.Params.ContentChanges[0].Text)
		if err != nil {
			return fmt.Errorf("failed to get urls and ignores: %w", err)
		}
		s.URLs[request.Params.TextDocument.URI] = append(s.URLs[request.Params.TextDocument.URI], data.URLs...)
		diags, err := s.GetDiagnosticsForFile(&request.Params.ContentChanges[0].Text, data)
		if err != nil {
			return fmt.Errorf("failed to get diagnostics for file: %w", err)
		}
		response.Params.Diagnostics = diags
		return nil
	})
	if err := eg.Wait(); err != nil {
		return response, fmt.Errorf("failed to get urls and ignores: %w", err)
	}
	return response, nil
}
