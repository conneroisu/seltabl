package analysis

import (
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
)

// Hover returns a hover response for the given uri and position
func (s *State) Hover(
	id int,
	uri string,
	position lsp.Position,
) (*lsp.HoverResponse, error) {
	text := s.Documents[uri]
	return &lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf(
				"Position: %s, document: %s",
				position.String(),
				text,
			),
			/*  */},
	}, nil
}
