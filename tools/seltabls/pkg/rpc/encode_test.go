package rpc_test

import (
	"context"
	"testing"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/stretchr/testify/assert"
	"go.lsp.dev/protocol"
)

// var (
//         EncodingInterfaces = []interface{}{
//                 lsp.CompletionResponse{
//                         Result: []lsp.CompletionItem{
//                                 {
//                                         Label:         "Test",
//                                         Detail:        "Test",
//                                         Documentation: "Test",
//                                         Kind:          lsp.CompletionKindMethod,
//                                 },
//                         },
//                 },
//                 lsp.HoverResponse{},
//         }
// )

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
	ctx := context.Background()
	expected := "Content-Length: 101\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":[{\"detail\":\"Test\",\"documentation\":\"Test\",\"kind\":2,\"label\":\"Test\"}]}\n"
	actual, err := rpc.Encode(ctx,
		lsp.TextDocumentCompletionResponse{
			Response: lsp.Response{
				RPC: lsp.RPCVersion,
				ID:  1,
			},
			Result: []protocol.CompletionItem{
				{
					Label:         "Test",
					Detail:        "Test",
					Documentation: "Test",
					Kind:          protocol.CompletionItemKindMethod,
				},
			},
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
