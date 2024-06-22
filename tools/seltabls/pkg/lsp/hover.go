package lsp

// HoverRequest is sent from the client to the server to request hover
// information.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
type HoverRequest struct {
	// HoverRequest embeeds the request struct.
	Request
	// Params are the parameters for the hover request.
	Params HoverParams `json:"params"`
}

// HoverParams is the parameters for a hover request.
type HoverParams struct {
	// TextDocumentPositionParams is the text document position parameters.
	TextDocumentPositionParams
}

// HoverResponse is the response from the server to a hover request.
type HoverResponse struct {
	// Response is the response for the hover request.
	Response
	// Result is the result for the hover request.
	Result HoverResult `json:"result"`
}

// HoverResult is a result from a hover request to the client from the
// language server.
type HoverResult struct {
	// Contents is the contents for the hover result.
	Contents string `json:"contents"`
}
