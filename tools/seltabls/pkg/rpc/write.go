package rpc

import (
	"context"
	"fmt"
	"io"
	"sync"
)

var (
	mu sync.Mutex
)

// WriteResponse writes a message to the writer
func WriteResponse(
	ctx context.Context,
	writer *io.Writer,
	msg MethodActor,
) error {
	mu.Lock()
	defer mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			reply, err := Encode(
				ctx,
				msg,
			)
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
		}
	}
}
