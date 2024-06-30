package rpc

import (
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
		return "", fmt.Errorf("failed to marshal message w/ method: %s: %w", msg.Method(), err)
	}
	return fmt.Sprintf(
		"Content-Length: %d\r\n\r\n%s",
		len(content),
		content,
	), nil
}
