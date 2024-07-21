package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// BaseMessage is the base message for a rpc message
type BaseMessage struct {
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Content []byte `json:"-"`
	Header  string `json:"-"`
}

// DecodeMessage decodes a rpc message
// returns the method, content, and error
func DecodeMessage(msg []byte) (*BaseMessage, error) {
	// Split the message into header and content
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return nil, fmt.Errorf("no header found")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse content length %s: %w",
			string(contentLengthBytes),
			err,
		)
	}
	var baseMessage BaseMessage
	err = json.Unmarshal(content[:contentLength], &baseMessage)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal base message: %w",
			err,
		)
	}
	baseMessage.Content = content[:contentLength]
	baseMessage.Header = string(header)
	return &baseMessage, nil
}
