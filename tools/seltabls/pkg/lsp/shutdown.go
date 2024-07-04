package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// ShutdownRequest is the request
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
type ShutdownRequest struct {
	Request
}

// Method returns the method for the shutdown request
func (r ShutdownRequest) Method() methods.Method {
	return methods.MethodShutdown
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
func (r ShutdownResponse) Method() string {
	return "shutdown"
}

// NewShutdownResponse creates a new shutdown response
func NewShutdownResponse(request ShutdownRequest, err error) ShutdownResponse {
	return ShutdownResponse{
		Response: Response{
			RPC: "2.0",
			ID:  request.ID,
		},
		Result: nil,
		Error:  &err,
	}
}

// ExitRequest is a struct for the exit request
type ExitRequest struct {
	Request
}

// Method returns the method for the exit request
func (r ExitRequest) Method() methods.Method {
	return methods.MethodExit
}
