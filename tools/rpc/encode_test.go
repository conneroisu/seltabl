package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncode tests the EncodeMessage function
func TestEncode(t *testing.T) {
	expected := "Content-Length: 66\r\n\r\n{\"Testing\":true,\"Method\":\"EncodeMessage\",\"Params\":\"Hello, World!\"}"
	actual := EncodeMessage(
		EncodingExample{
			Testing: true,
			Method:  "EncodeMessage",
			Params:  "Hello, World!",
		},
	)
	assert.Equal(t, expected, actual)
}
