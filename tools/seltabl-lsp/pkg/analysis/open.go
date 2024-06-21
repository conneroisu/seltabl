package analysis

import (
	"context"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/parsers"
)

// OpenDocument opens a document in the state and returns any diagnostics for the document
//
// uri is the uri of the document
//
// content is the content of the document
func (s *State) OpenDocument(
	uri string,
	content *string,
) (diags []lsp.Diagnostic, err error) {
	ctx := context.Background()
	s.Documents[uri] = *content
	data, err := parsers.ParseStructComments(*content)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
		return nil, err
	}
	for _, url := range data.URLs {
		s.Selectors[uri], err = s.getSelectors(
			ctx,
			[]string{url},
			data.IgnoreElements,
		)
		if err != nil {
			s.Logger.Printf("failed to get selectors: %s\n", err)
		}
	}
	diags, err = s.GetDiagnosticsForFile(content, data)
	if err != nil {
		s.Logger.Printf("failed to get diagnostics for file: %s\n", err)
	}
	return diags, nil
}
