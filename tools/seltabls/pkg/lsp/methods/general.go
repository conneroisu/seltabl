package methods

import (
	"context"
	"fmt"
	"io"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
)

// Notification Methods
const (
	// MethodInitialize is the initialize notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
	MethodInitialize Method = "initialize"

	// MethodInitialized is the initialized notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialized
	MethodInitialized Method = "initialized"

	// MethodShutdown is the shutdown notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
	MethodShutdown Method = "shutdown"

	// MethodExit is the exit notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#exit
	MethodExit Method = "exit"

	// MethodCancelRequest is the cancel request method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#cancelRequest
	MethodCancelRequest Method = "$/cancelRequest"
)

// MethodHandler is a function for handling a method
type MethodHandler func(ctx context.Context, writer *io.Writer, state *analysis.State, msg []byte) error

// MethodHandlers is a map of method handlers
type MethodHandlers map[Method]MethodHandler

// Handle handles a message
func (m MethodHandlers) Handle(ctx context.Context, writer *io.Writer, state *analysis.State, msg []byte) error {
	handler, ok := m[Method(msg[0])]
	if !ok {
		return fmt.Errorf("no handler found for method: %s", Method(msg[0]))
	}
	return handler(ctx, writer, state, msg)
}
