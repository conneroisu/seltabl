package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// BaseMessage is the base message for a rpc message
type BaseMessage struct {
	Testing bool   `json:"testing"`
	ID      *int   `json:"id,omitempty"`
	Method  string `json:"method"`
	Content []byte `json:"-"`
}

// DecodeMessage decodes a rpc message
func DecodeMessage(msg []byte) (BaseMessage, error) {
	// Split the message into header and content
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return BaseMessage{}, fmt.Errorf("no header found")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return BaseMessage{}, fmt.Errorf("failed to parse content length: %w", err)
	}
	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return BaseMessage{}, fmt.Errorf("failed to unmarshal base message: %w", err)
	}
	// return baseMessage.Method, content[:contentLength], nil
	baseMessage.Content = content[:contentLength]
	return baseMessage, nil
}
