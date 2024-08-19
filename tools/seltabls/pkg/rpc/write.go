package rpc

import (
	"context"
	"fmt"
	"io"
	"sync"
)

// Writer is a struct for writing a response to a writer.
//
// It implements the io.Writer interface.
type Writer struct {
	io.Writer
	mu sync.Mutex
}

// NewWriter creates a new writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w}
}

// WriteResponse writes a message to the writer
func (w *Writer) WriteResponse(
	ctx context.Context,
	msg MethodActor,
) error {
	w.mu.Lock()
	defer w.mu.Unlock()
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
			res, err := w.Write([]byte(reply))
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
