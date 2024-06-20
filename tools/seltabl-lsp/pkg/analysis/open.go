package analysis

import (
	"context"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/parsers"
)

// OpenDocument opens a document in the state and returns any diagnostics for the document
func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	ctx := context.Background()
	s.Documents[uri] = text
	out, err := parsers.ParseStructComments(text)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
		return nil
	}
	for _, url := range out.URLs {
		s.Selectors[uri], err = s.getSelectors(
			ctx,
			[]string{url},
			out.IgnoreElements,
		)
		if err != nil {
			s.Logger.Printf("failed to get selectors: %s\n", err)
		}
	}
	return getDiagnosticsForFile(text)
}
