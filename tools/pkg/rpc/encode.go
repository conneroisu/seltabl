package rpc

import (
	"encoding/json"
	"fmt"
)

// EncodeMessage encodes a message into a string
func EncodeMessage(msg any) (string, error) {
	content, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content), nil
}
