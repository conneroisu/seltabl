package analysis

import "github.com/conneroisu/seltabl/tools/pkg/lsp"

// Definition returns the definition.
func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// In real life, this would look up the definition

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}
