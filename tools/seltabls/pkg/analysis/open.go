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
	uri string,
	content *string,
) (diags []lsp.Diagnostic, err error) {
	eg, ctx := errgroup.WithContext(ctx)
	s.Documents[uri] = *content
	data, err := parsers.ParseStructComments(*content)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
		return diags, nil
	}
	for _, url := range data.URLs {
		eg.Go(func() error {
			s.Selectors[uri], err = s.GetSelectors(
				ctx,
				url,
				data.IgnoreElements,
			)
			if err != nil {
				return fmt.Errorf("failed to get selectors for url (%s): %w", url, err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return diags, fmt.Errorf("failed to get selectors for urls: %w", err)
	}
	diags, err = s.GetDiagnosticsForFile(content, data)
	if err != nil {
		s.Logger.Printf("failed to get diagnostics for file: %s\n", err)
	}
	return diags, nil
}
