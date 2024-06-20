package analysis

import (
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/parsers"
)

// Hover returns a hover response for the given uri and position
func (s *State) Hover(
	id int,
	uri string,
	position lsp.Position,
) (*lsp.HoverResponse, error) {
	document := s.Documents[uri]
	check, err := s.CheckPosition(position, document)
	if err != nil {
		return nil, fmt.Errorf(
			"analysis.Hover() error: %v",
			err,
		)
	}
	switch check {
	case parsers.StateInTag:
		s.Logger.Println("Found position in struct tag")
	case parsers.StateInTagValue:
		s.Logger.Println("Found position in struct tag value")
	case parsers.StateAfterColon:
		s.Logger.Println("Found position in struct tag after colon")
	case parsers.StateInvalid:
	default:
		return nil, nil
	}
	return &lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		// Result: lsp.HoverResult{
		//         Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
		// },
	}, nil
}
