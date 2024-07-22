package rpc

import (
	"bytes"
	"fmt"
	"strconv"
)

// Split splits a byte slice into a header and content.
//
// It returns the advance, token, and error.
func Split(data []byte, _ bool) (int, []byte, error) {
	var err error
	var advance int
	var token []byte
	var header, content []byte
	var found bool
	var contentLengthBytes []byte
	var contentLength int
	header, content, found = bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}
	// Content-Length: <number>
	contentLengthBytes = header[len("Content-Length: "):]
	contentLength, err = strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse content length: %w", err)
	}
	if len(content) < contentLength {
		return 0, nil, nil
	}
	advance = len(header) + 4 + contentLength
	token = data[:advance]
	return advance, token, nil
}
