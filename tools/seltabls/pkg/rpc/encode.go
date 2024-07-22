package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
)

// MethodActor is a type for responding to a request
type MethodActor interface {
	Method() methods.Method
}

// Encode encodes a message into a string
//
// It uses the json library to encode the message
// and returns a string representation of the encoded message with
// a Content-Length header.
//
// It also returns an error if there is an error while encoding the message.
func Encode(
	ctx context.Context,
	msg MethodActor,
) (string, error) {
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		var err error
		var buffer *bytes.Buffer
		var encoder *json.Encoder
		buffer = &bytes.Buffer{}
		encoder = json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		err = encoder.Encode(msg)
		if err != nil {
			return "", err
		}
		var body []byte
		body = bytes.TrimSuffix(buffer.Bytes(), []byte("\n"))
		log.Debugf(
			"wrote msg [%d] (%s): %s",
			len(body),
			msg.Method(),
			string(body),
		)
		result := fmt.Sprintf(
			"Content-Length: %d\r\n\r\n%s",
			len(body),
			string(body),
		)
		return result, nil
	}
}
