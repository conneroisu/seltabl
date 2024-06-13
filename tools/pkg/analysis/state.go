package analysis

import (
	"github.com/conneroisu/seltabl/tools/pkg/lsp"
)

// State is the state of the document analysis
type State struct {
	// Map of file names to contents
	Documents map[string]string
}

// NewState returns a new state with no documents
func NewState() State {
	return State{Documents: map[string]string{}}
}

// LineRange returns a range of a line in a document
func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
