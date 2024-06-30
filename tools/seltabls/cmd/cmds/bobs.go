package cmds

import (
	"context"
	"fmt"
	"io"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"golang.org/x/sync/errgroup"
)

// WriteResponse writes a message to the writer
func WriteResponse(
	ctx context.Context,
	state *analysis.State,
	writer *io.Writer,
	msg rpc.MethodActor,
) error {
	responseCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(responseCtx)
	eg.Go(func() error {
		reply, err := rpc.EncodeMessage(msg)
		if err != nil {
			return fmt.Errorf(
				"failed to encode response to request (%s): %w",
				msg.Method(),
				err,
			)
		}
		state.Logger.Printf("sending response (%s): %s", msg.Method(), string(reply))
		res, err := (*writer).Write([]byte(reply))
		if err != nil {
			return fmt.Errorf(
				"failed to encode response to request (%s): %w",
				msg.Method(),
				err,
			)
		}
		if res != len(reply) {
			return fmt.Errorf(
				"failed to write all of response to request (%s): %w",
				msg.Method(),
				err,
			)
		}
		return nil
	})
	return eg.Wait()
}
