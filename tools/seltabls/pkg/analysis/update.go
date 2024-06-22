package analysis

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
)

// UpdateDocument updates the state with the given document
func (s *State) UpdateDocument(
	uri, content string,
) (diags []lsp.Diagnostic, err error) {
	s.Documents[uri] = content
	data, err := parsers.ParseStructComments(content)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
		return nil, err
	}
	diags, err = s.GetDiagnosticsForFile(&content, data)
	if err != nil {
		s.Logger.Printf("failed to get diagnostics for file: %s\n", err)
	}
	return diags, nil
}
