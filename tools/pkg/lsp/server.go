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
