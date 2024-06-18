package lsp

// ShutdownRequest is the request
type ShutdownRequest struct {
	ID *string `json:"id,omitempty"`
}

// ShutdownResponse is the response to a ShutdownRequest.
type ShutdownResponse struct {
	Response
	Error *error `json:"error,omitempty"`
}
