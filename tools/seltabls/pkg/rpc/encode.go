package rpc

import (
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
	content, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Content-Length: %d\r\n\r\n%s",
		len(content),
		content,
	), nil
}
