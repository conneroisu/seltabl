package analysis

import (
	"github.com/conneroisu/seltabl/tools/pkg/lsp"
)

// Hover returns a hover response for the given uri and position
func (s *State) Hover(
	id int,
	uri string,
	position lsp.Position,
) lsp.HoverResponse {
	document := s.Documents[uri]
	if !determineIfRespondHover(document, position) {
		return lsp.HoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: lsp.HoverResult{
				Contents: "",
			},
		}
	}

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		// Result: lsp.HoverResult{
		//         Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		// },
	}
}

// determineIfRespondHover determines to response to a hover request
func determineIfRespondHover(_ string, _ lsp.Position) bool {
	return true
}
