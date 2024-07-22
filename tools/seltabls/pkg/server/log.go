package server

import (
	"context"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"go.lsp.dev/protocol"
)

var (
	writer io.Writer = os.Stdout
)

// LogMessage sends a log message to the client.
func LogMessage(
	ctx context.Context,
	msg string,
	typ protocol.MessageType,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var err error
			err = rpc.WriteResponse(ctx, &writer, lsp.LogMessageNotification{
				Notification: lsp.Notification{
					RPC:    lsp.RPCVersion,
					Method: "window/logMessage",
				},
				Params: protocol.LogMessageParams{
					Message: msg,
					Type:    typ,
				},
			})
			if err != nil {
				log.Errorf("failed to write log message: %s", err)
			}
			return
		}
	}

}
