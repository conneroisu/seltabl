package rpc

import (
	"encoding/json"
	"fmt"
)

// EncodeMessage encodes a message into a string
func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}
