package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// CancelRequest is sent from the client to the server to cancel a request.
type CancelRequest struct {
	// CancelRequest embeds the Request struct
	Request
	// ID is the id of the request to be cancelled.
	ID string `json:"id"`
	// Params are the parameters for the request to be cancelled.
	Params CancelParams `json:"params"`
}

// Method returns the method for the cancel request
func (r CancelRequest) Method() methods.Method {
	return methods.MethodCancelRequest
}

// CancelParams are the parameters for a cancel request.
type CancelParams struct {
	// ID is the id of the request to be cancelled.
	ID string `json:"id"`
}

// CancelResponse is the response for a cancel request.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_cancel
type CancelResponse struct {
	// CancelResponse embeds the Response struct
	Response
}

// Method returns the method for the cancel response
func (r *CancelResponse) Method() methods.Method {
	return methods.MethodCancelRequest
}
