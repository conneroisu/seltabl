package analysis

import "github.com/conneroisu/seltabl/tools/pkg/lsp"

// OpenDocument opens a document in the state and returns any diagnostics for the document
func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}
