package cmds

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/spf13/cobra"
)

// LSPHandler is a struct for the LSP server
type LSPHandler func(ctx context.Context, writer *io.Writer, state *analysis.State, msg rpc.BaseMessage) error

// handleCtx is a struct for the handle context.
type handleCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewLSPCmd creates a new command for the lsp subcommand
func NewLSPCmd(
	ctx context.Context,
	writer io.Writer,
	reader io.Reader,
	handle LSPHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp", // the name of the command
		Short: "A command line tooling for package that parsing html tables and elements into structs",
		Long: `
CLI and Language Server for the seltabl package.

Language server provides completions, hovers, and code actions for seltabl defined structs.
	
CLI provides a command line tool for verifying, linting, and reporting on seltabl defined structs.
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			scanner := bufio.NewScanner(reader)
			scanner.Split(rpc.Split)
			state, err := analysis.NewState()
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			_, cancel := context.WithCancel(ctx)
			defer cancel()
			ctxs := make(map[int]handleCtx)
			for scanner.Scan() {
				hCtx, cancel := context.WithCancel(ctx)
				defer cancel()
				decoded, err := rpc.DecodeMessage(scanner.Bytes())
				if err != nil {
					return fmt.Errorf("failed to decode message: %w", err)
				}
				if decoded.Method == string(methods.MethodCancelRequest) {
					ctxs[decoded.ID].cancel()
					delete(ctxs, decoded.ID)
					continue
				}
				ctxs[decoded.ID] = handleCtx{ctx: hCtx, cancel: cancel}
				err = handle(hCtx, &writer, &state, decoded)
				if err != nil {
					state.Logger.Printf("failed to handle message: %s\n", err)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scanner error: %w", err)
			}
			return nil
		},
	}
	return cmd
}
