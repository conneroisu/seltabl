package rpc_test

import (
	"testing"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/stretchr/testify/assert"
)

// EncodingExample is a test struct
type EncodingExample struct {
	// Testing is a test field
	Testing bool `json:"testing"`
	// Method is the method for the request
	Meth string `json:"method"`
	// Params are the parameters for the request
	Params string `json:"params"`
}

// Method returns the method for the request
func (e EncodingExample) Method() string {
	return e.Meth
}

// TestEncode tests the EncodeMessage function
func TestEncode(t *testing.T) {
	expected := "Content-Length: 100\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":[{\"label\":\"Test\",\"detail\":\"Test\",\"documentation\":\"Test\",\"kind\":2}]}"
	actual, err := rpc.EncodeMessage(
		lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  1,
			},
			Result: []lsp.CompletionItem{
				{
					Label:         "Test",
					Detail:        "Test",
					Documentation: "Test",
					Kind:          lsp.Method,
				},
			},
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
