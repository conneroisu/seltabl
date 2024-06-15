package analysis

import "github.com/conneroisu/seltabl/tools/pkg/lsp"

// OpenDocument opens a document in the state and returns any diagnostics for the document
func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	urls, ignores, err := s.getUrlsAndIgnores(text)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
		return nil
	}
	for _, url := range urls {
		s.Selectors[uri], err = s.getSelectors(url, ignores)
		if err != nil {
			s.Logger.Printf("failed to get selectors: %s\n", err)
		}
	}
	return getDiagnosticsForFile(text)
}
