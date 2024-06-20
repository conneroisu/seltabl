package lsp

import (
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

// Server is a struct for the LSP server
type Server interface {
	// bun.QueryHook is an embedded interface within the Server interface
	bun.QueryHook
	// HandleMessage handles a message from the client
	HandleMessage(msg []byte) error
	// ReturnCmd returns the command for starting the LSP server
	ReturnCmd() *cobra.Command
}

// LineRange returns a range of a line in a document
//
// line is the line number
//
// start is the start character of the range
//
// end is the end character of the range
func LineRange(line, start, end int) Range {
	return Range{
		Start: Position{
			Line:      line,
			Character: start,
		},
		End: Position{
			Line:      line,
			Character: end,
		},
	}
}
