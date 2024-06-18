package analysis

import (
	"context"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
)

// OpenDocument opens a document in the state and returns any diagnostics for the document
func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	ctx := context.Background()
	s.Documents[uri] = text
	urls, ignores, err := s.getUrlsAndIgnores(text)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
		return nil
	}
	s.Logger.Printf("urls: %v\n", urls)
	s.Logger.Printf("ignores: %v\n", ignores)
	for _, url := range urls {
		s.Selectors[uri], err = s.getSelectors(ctx, []string{url}, ignores)
		if err != nil {
			s.Logger.Printf("failed to get selectors: %s\n", err)
		}
	}
	return getDiagnosticsForFile(text)
}
