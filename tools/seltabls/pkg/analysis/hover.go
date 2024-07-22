package analysis

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/http"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/yosssi/gohtml"
	"go.lsp.dev/protocol"
)

// NewHoverResponse returns a hover response for the given uri and position
func NewHoverResponse(
	req lsp.HoverRequest,
	s *State,
) (response *lsp.HoverResponse, err error) {
	response = &lsp.HoverResponse{
		Response: lsp.Response{
			RPC: lsp.RPCVersion,
			ID:  req.ID,
		},
		Result: lsp.HoverResult{},
	}
	text := s.Documents[req.Params.TextDocument.URI]
	urls := s.URLs[req.Params.TextDocument.URI]
	doc, err := http.DefaultClientGet(urls[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get the content of the url: %w", err)
	}
	res, err := s.GetSelectorHover(req.Params.Position, text, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to get hover: %w", err)
	}
	response.Result = res
	return response, nil
}

// GetSelectorHover checks if the position is within the struct tag
func (s *State) GetSelectorHover(
	position protocol.Position,
	text string,
	doc *goquery.Document,
) (res lsp.HoverResult, err error) {
	var inValue bool
	// Create a new token file set
	fset := token.NewFileSet()
	position.Line = position.Line + 1
	// Parse the source code from a new buffer
	node, err := parser.ParseFile(
		fset,
		"",
		bytes.NewBufferString(text),
		parser.Trace,
	)
	if err != nil {
		return res, fmt.Errorf("failed to parse struct: %w", err)
	}
	// Find the struct node in the AST
	structNodes := parsers.FindStructNodes(node)
	for i := range structNodes {
		// Check if the position is within the struct node
		inPosition := parsers.IsPositionInNode(structNodes[i], position, fset)
		// Check if the position is within a struct tag
		inTag := parsers.IsPositionInTag(structNodes[i], position, fset)
		if inPosition && inTag {
			var val string
			// Check if the position is within a struct tag value (i.e. value inside and including " and " characters)
			val, inValue = parsers.PositionInStructTagValue(
				structNodes[i],
				position,
				fset,
			)
			if !inValue {
				if parsers.PositionBeforeValue(position, text) != ':' &&
					parsers.PositionBeforeValue(position, text) != '"' {
					continue
				}
			}
			HTML, err := doc.Find(val).Parent().Html()
			if err != nil {
				return res, fmt.Errorf("failed to get html: %w", err)
			}
			HTML = gohtml.Format(HTML)
			res.Contents = HTML
			return res, nil
		}
	}
	return res, nil
}
