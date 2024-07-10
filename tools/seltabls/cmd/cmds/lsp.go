package cmds

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/server"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// LSPHandler is a struct for the LSP server
type LSPHandler func(ctx context.Context, state *analysis.State, msg rpc.BaseMessage) (rpc.MethodActor, error)

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
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			var scanner *bufio.Scanner
			var state *analysis.State
			scanner = bufio.NewScanner(reader)
			state, err = analysis.NewState()
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			mw := io.MultiWriter(state.Logger.Writer(), os.Stderr)
			state.Logger.SetOutput(mw)
			cmd.SetErr(mw)
			ctx, cancel := context.WithCancel(ctx)
			eg, _ := errgroup.WithContext(ctx)
			defer cancel()
			ctxs := make(map[int]handleCtx)
			scanner.Split(rpc.Split)
			for scanner.Scan() {
				eg.Go(func() error {
					input := scanner.Bytes()
					hCtx, cancel := context.WithCancel(context.Background())
					defer cancel()
					decoded, err := rpc.DecodeMessage(input)
					if err != nil {
						return fmt.Errorf("failed to decode message: %w", err)
					}
					handie := handleCtx{ctx: hCtx, cancel: cancel}
					content, err := json.MarshalIndent(&decoded.Content, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal content: %w", err)
					}
					state.Logger.Printf(
						"Received message: %s\n",
						content,
					)
					ctxs[decoded.ID] = handie
					if decoded.Method == string(methods.MethodCancelRequest) {
						ctxs[decoded.ID].cancel()
					}
					resp, err := handle(hCtx, state, decoded)
					if err != nil || resp == nil {
						state.Logger.Printf(
							`failed to handle message: %v
input: %s
id: %d
method: %s
ctx-canceled: %v
current-err: %v`,
							err,
							input,
							decoded.ID,
							decoded.Method,
							ctx.Err(),
							err,
						)
						return nil
					}
					err = server.WriteResponse(hCtx, &writer, resp)
					if err != nil {
						state.Logger.Printf(
							"failed to write response: %s\n",
							err,
						)
					}
					state.Logger.Printf("Sent response: %s", resp)
					return nil
				})
			}
			err = scanner.Err()
			if err != nil {
				return fmt.Errorf("scanner error: %w", err)
			}
			err = eg.Wait()
			if err != nil {
				return fmt.Errorf("failed to wait for context: %w", err)
			}
			return nil
		},
	}
	return cmd
}
