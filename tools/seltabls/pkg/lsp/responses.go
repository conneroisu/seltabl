package lsp

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// CancelResponse is the response for a cancel request.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_cancel
type CancelResponse struct {
	RPC string `json:"jsonrpc"`
	ID  int    `json:"id"`
}

// Method returns the method for the cancel response
func (r CancelResponse) Method() methods.Method {
	return methods.MethodCancelRequest
}

// TextDocumentCodeActionResponse is the response for a code action request.
type TextDocumentCodeActionResponse struct {
	// TextDocumentCodeActionResponse embeds the Response struct
	Response
	// Result is the result for the code action request.
	Result []protocol.CodeAction `json:"result"`
}

// Method returns the method for the code action response
func (r TextDocumentCodeActionResponse) Method() methods.Method {
	return methods.MethodRequestTextDocumentCodeAction
}

// HoverResponse is the response from the server to a hover request.
type HoverResponse struct {
	// Response is the response for the hover request.
	Response
	// Result is the result for the hover request.
	Result HoverResult `json:"result"`
}

// Method returns the method for the hover response
func (r HoverResponse) Method() methods.Method {
	return methods.MethodRequestTextDocumentHover
}

// HoverResult is a result from a hover request to the client from the
// language server.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
type HoverResult struct {
	// Contents is the contents for the hover result.
	Contents string `json:"contents"`
}

// InitializeResponse is a struct for the initialize response.
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
//
// It embeds the Response struct.
type InitializeResponse struct {
	Response
	// Result is the result of the initialize request
	Result protocol.InitializeResult `json:"result"`
}

// Method returns the method for the initialize response
func (r InitializeResponse) Method() methods.Method {
	return methods.MethodInitialize
}

// NewInitializeResponse creates a new initialize response.
func NewInitializeResponse(
	ctx context.Context,
	request *InitializeRequest,
) (*InitializeResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			return &InitializeResponse{
				Response: Response{
					RPC: RPCVersion,
					ID:  request.ID,
				},
				Result: protocol.InitializeResult{
					Capabilities: protocol.ServerCapabilities{
						TextDocumentSync: protocol.TextDocumentSyncOptions{
							OpenClose: true,
							Change:    protocol.TextDocumentSyncKindFull,
						},
						CompletionProvider: &protocol.CompletionOptions{
							ResolveProvider: false,
							TriggerCharacters: []string{
								":", "\"",
							},
						},
						HoverProvider:      true,
						DefinitionProvider: false,
						CodeActionProvider: false,
					},
					ServerInfo: &protocol.ServerInfo{
						Name:    "seltabl_lsp",
						Version: "0.0.0.5.0.0-beta1.final",
					},
				},
			}, nil
		}
	}
}

// ShutdownResponse is the response to a ShutdownRequest.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
type ShutdownResponse struct {
	Response
	Result *bool  `json:"result"`
	Error  *error `json:"error,omitempty"`
}

// Method returns the method for the shutdown response
func (r ShutdownResponse) Method() methods.Method {
	return methods.MethodShutdown
}

// NewShutdownResponse creates a new shutdown response
func NewShutdownResponse(
	request ShutdownRequest,
	err error,
) (ShutdownResponse, error) {
	return ShutdownResponse{
		Response: Response{
			RPC: RPCVersion,
			ID:  request.ID,
		},
		Result: nil,
		Error:  &err,
	}, nil
}
