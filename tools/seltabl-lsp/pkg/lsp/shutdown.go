package lsp

// ShutdownRequest is the request
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
type ShutdownRequest struct {
	Request
}

// ShutdownResponse is the response to a ShutdownRequest.
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
type ShutdownResponse struct {
	Response
	Error *error `json:"error,omitempty"`
}
