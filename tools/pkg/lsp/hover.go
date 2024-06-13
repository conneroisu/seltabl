package lsp

// HoverRequest is sent from the client to the server to request hover
// information.
type HoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}

// HoverParams is the parameters for a hover request.
type HoverParams struct {
	TextDocumentPositionParams
}

// HoverResponse is the response from the server to a hover request.
type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}

// HoverResult is a result from a hover request to the client from the
// language server.
type HoverResult struct {
	Contents string `json:"contents"`
}
