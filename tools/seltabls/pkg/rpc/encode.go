package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MethodActor is a type for responding to a request
type MethodActor interface {
	Method() string
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
	var buf bytes.Buffer
	json.HTMLEscape(&buf, content)
	return fmt.Sprintf(
		"Content-Length: %d\r\n\r\n%s",
		buf.Len(),
		buf.String(),
	), nil
}
