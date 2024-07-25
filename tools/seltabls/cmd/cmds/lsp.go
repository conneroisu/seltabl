package cmds

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"github.com/spf13/cobra"
	"go.lsp.dev/uri"
)

// LSPHandler is a struct for the LSP server
type LSPHandler func(
	ctx context.Context,
	cancel *context.CancelFunc, // cancel is a pointer to the cancel function to avoid copying
	msg *rpc.BaseMessage, // msg is a pointer to the message to avoid copying
	db *data.Database[master.Queries], // db is a pointer to the database to avoid copying
	documents *safe.Map[uri.URI, string],
	selectors *safe.Map[uri.URI, []master.Selector],
	urls *safe.Map[uri.URI, []string],
) (rpc.MethodActor, error)

// NewLSPCmd creates a new command for the lsp subcommand
func NewLSPCmd(
	ctx context.Context,
	writer io.Writer,
	reader io.Reader,
	handle LSPHandler,
	db *data.Database[master.Queries],
) *cobra.Command {
	return &cobra.Command{
		Use:   "lsp", // the name of the command
		Short: "A command line tooling for package that parsing html tables and elements into structs",
		Long: `
CLI and Language Server for the seltabl package.

Language server provides completions, hovers, and code actions for seltabl defined structs.
	
CLI provides a command line tool for verifying, linting, and reporting on seltabl defined structs.
`,
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			documents := safe.NewSafeMap[uri.URI, string]()
			selectors := safe.NewSafeMap[uri.URI, []master.Selector]()
			urls := safe.NewSafeMap[uri.URI, []string]()
			scanner := bufio.NewScanner(reader)
			scanner.Split(rpc.Split)
			lspCtx, lspCancel := context.WithCancel(ctx)
			defer lspCancel()
			for scanner.Scan() {
				go func() {
					decoded, err := rpc.DecodeMessage(scanner.Bytes())
					if err != nil {
						log.Errorf("failed to decode message: %s", err)
					}
					hCtx, cancel := context.WithCancel(lspCtx)
					lsp.CancelMap.Set(decoded.ID, cancel)
					defer lsp.CancelMap.Delete(decoded.ID)
					log.Debugf(
						"received message [%s] (%s): %s",
						decoded.Header,
						decoded.Method,
						decoded.Content,
					)
					resp, err := handle(
						hCtx,
						&lspCancel,
						decoded,
						db,
						documents,
						selectors,
						urls,
					)
					if err != nil {
						log.Errorf(
							"failed to handle message (%s): %s",
							decoded.Method,
							err,
						)
					}
					if !isNull(resp) {
						err = rpc.WriteResponse(hCtx, &writer, resp)
						if err != nil {
							log.Errorf(
								"failed to write (%s) response: %s",
								resp.Method(),
								err,
							)
						}
					}
				}()
			}
			err = scanner.Err()
			if err != nil {
				return fmt.Errorf("scanner error: %w", err)
			}
			return nil
		},
	}
}

// isNull checks if the given interface is nil or points to a nil value
func isNull(i interface{}) bool {
	if i == nil {
		return true
	}
	// Use reflect.ValueOf only if the kind is valid for checking nil
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return v.IsNil()
	}
	return false
}
