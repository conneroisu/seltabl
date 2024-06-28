package rpc

import (
	"testing"
)

// TestDecode tests the decode function
func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	message, err := DecodeMessage([]byte(incomingMessage))
	contentLength := len(message.Content)
	if err != nil {
		t.Fatal(err)
	}

	if contentLength != 15 {
		t.Fatalf("Expected: 16, Got: %d", contentLength)
	}

	if message.Method != "hi" {
		t.Fatalf("Expected: 'hi', Got: %s", message.Method)
	}
}
