package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/charmbracelet/log"
)

// WriteResponse writes a message to the writer
func WriteResponse(
	ctx context.Context,
	writer *io.Writer,
	msg MethodActor,
) error {
	go func() {
		log.Debugf(
			"sent message [%d] (%s): %s",
			len(marshal(msg)),
			msg.Method(),
			marshal(msg),
		)
	}()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			reply, err := EncodeMessage(msg)
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

func marshal(mA MethodActor) string {
	b, err := json.Marshal(mA)
	if err != nil {
		return ""
	}
	return string(b)
}
