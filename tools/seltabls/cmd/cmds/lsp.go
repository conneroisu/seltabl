package cmds

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/server"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// LSPHandler is a struct for the LSP server
type LSPHandler func(ctx context.Context, state *analysis.State, msg rpc.BaseMessage, cancel context.CancelFunc) (rpc.MethodActor, error)

var (
	mu    sync.Mutex
	cnMap = make(map[int]context.CancelFunc)
)

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
		RunE: func(cmd *cobra.Command, _ []string) error {
			scanner := bufio.NewScanner(reader)
			scanner.Split(rpc.Split)
			state, err := analysis.NewState()
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			cmd.SetErr(state.Logger.Writer())
			lspCtx, lspCancel := context.WithCancel(ctx)
			defer lspCancel()
			eg, hCtx := errgroup.WithContext(lspCtx)
			for scanner.Scan() {
				eg.Go(func() error {
					decoded, err := rpc.DecodeMessage(scanner.Bytes())
					if err != nil {
						return fmt.Errorf("failed to decode message: %w", err)
					}
					hCtx, cancel := context.WithCancel(hCtx)
					mu.Lock()
					cnMap[decoded.ID] = cancel
					mu.Unlock()
					if decoded.Method == string(methods.MethodCancelRequest) {
						log.Debugf("canceling request: %s", decoded.Method)
						cnMap[decoded.ID]()
						mu.Lock()
						delete(cnMap, decoded.ID)
						mu.Unlock()
						return nil
					}
					log.Debugf("received message: %s", decoded.Method)
					resp, err := handle(hCtx, &state, *decoded, lspCancel)
					if err != nil {
						log.Errorf(
							"failed to handle message (%s): %s",
							decoded.Method,
							err,
						)
						return nil
					}
					if isNull(resp) || resp == nil {
						return nil
					}
					err = server.WriteResponse(hCtx, &writer, resp)
					if err != nil {
						log.Errorf(
							"failed to write (%s) response: %s\n",
							resp.Method(),
							err,
						)
					}
					go log.Debugf("sent message: %s", marshal(resp))
					mu.Lock()
					delete(cnMap, decoded.ID)
					mu.Unlock()
					return nil
				})
			}
			err = scanner.Err()
			if err != nil {
				return fmt.Errorf("scanner error: %w", err)
			}
			return nil
		},
	}
	return cmd
}

func marshal(mA rpc.MethodActor) string {
	b, err := json.Marshal(mA)
	if err != nil {
		return ""
	}
	return string(b)
}

func isNull(i interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("recovered from panic: %v", r)
		}
	}()
	return i == nil || reflect.ValueOf(i).IsNil()
}
