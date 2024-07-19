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
func (r ShutdownResponse) Method() methods.Method {
	return methods.MethodShutdown
}

// NewShutdownResponse creates a new shutdown response
func NewShutdownResponse(request ShutdownRequest, err error) (ShutdownResponse, error) {
	return ShutdownResponse{
		Response: Response{
			RPC: RPCVersion,
			ID:  request.ID,
		},
		Result: nil,
		Error:  &err,
	}, nil
}

// ExitRequest is a struct for the exit request
type ExitRequest struct {
	Request
}

// Method returns the method for the exit request
func (r ExitRequest) Method() methods.Method {
	return methods.MethodNotificationExit
}
