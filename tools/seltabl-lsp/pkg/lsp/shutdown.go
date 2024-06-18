package lsp

// ShutdownRequest is the request
type ShutdownRequest struct {
	Request
}

// ShutdownResponse is the response to a ShutdownRequest.
type ShutdownResponse struct {
	Response
	Error *error `json:"error,omitempty"`
}
