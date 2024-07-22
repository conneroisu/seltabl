package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
)

// MethodActor is a type for responding to a request
type MethodActor interface {
	Method() methods.Method
}

// EncodeMessage encodes a message into a string
//
// It uses the json library to encode the message
// and returns a string representation of the encoded message with
// a Content-Length header.
//
// It also returns an error if there is an error while encoding the message.
func EncodeMessage(msg MethodActor) (string, error) {
	var err error
	var buffer *bytes.Buffer
	var encoder *json.Encoder
	buffer = &bytes.Buffer{}
	encoder = json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(msg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Content-Length: %d\r\n\r\n%s",
		buffer.Len(),
		buffer.String(),
	), nil
}
