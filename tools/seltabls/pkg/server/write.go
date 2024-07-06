package server

import (
	"context"
	"fmt"
	"io"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"golang.org/x/sync/errgroup"
)

// WriteResponse writes a message to the writer
func WriteResponse(
	ctx context.Context,
	writer *io.Writer,
	msg rpc.MethodActor,
) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		reply, err := rpc.EncodeMessage(msg)
		if err != nil {
			return fmt.Errorf(
				"failed to encode response to request (%s): %w",
				msg.Method(),
				err,
			)
		}
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
				"failed writing all of response to (%s) request: %w",
				msg.Method(),
				err,
			)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf(
			"failed to write message (%s): %w",
			msg.Method(),
			err,
		)
	}
	return nil
}
