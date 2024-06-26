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
	Method  string `json:"method"`
}

// DecodeMessage decodes a rpc message
func DecodeMessage(msg []byte) (string, []byte, error) {
	// Split the message into header and content
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, fmt.Errorf("no header found")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse content length: %w", err)
	}
	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal base message: %w", err)
	}
	return baseMessage.Method, content[:contentLength], nil
}
