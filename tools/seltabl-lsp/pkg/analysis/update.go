package analysis

import "github.com/conneroisu/seltabl/tools/pkg/lsp"

// UpdateDocument updates the state with the given document
func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}
