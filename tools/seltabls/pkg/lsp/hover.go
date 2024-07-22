package lsp

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// HoverRequest is sent from the client to the server to request hover
// information.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
type HoverRequest struct {
	// HoverRequest embeeds the request struct.
	Request
	// Params are the parameters for the hover request.
	Params protocol.HoverParams `json:"params"`
}

// Method returns the method for the hover request
func (r HoverRequest) Method() methods.Method {
	return methods.MethodRequestTextDocumentHover
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
