package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// EncodingExample is a test struct
type EncodingExample struct {
	// Testing is a test field
	Testing bool
	// Method is the method for the request
	Method string
	// Params are the parameters for the request
	Params string
}

// TestEncode tests the EncodeMessage function
func TestEncode(t *testing.T) {
	expected := "Content-Length: 66\r\n\r\n{\"Testing\":true,\"Method\":\"EncodeMessage\",\"Params\":\"Hello, World!\"}"
	actual, err := EncodeMessage(
		EncodingExample{
			Testing: true,
			Method:  "EncodeMessage",
			Params:  "Hello, World!",
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
