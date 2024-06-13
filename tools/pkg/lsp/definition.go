package lsp

// DefinitionRequest is sent from the client to the server to resolve a
// definition.
type DefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}

// DefinitionParams is the set of parameters that are sent in a
// DefinitionRequest.
type DefinitionParams struct {
	TextDocumentPositionParams
}

// DefinitionResponse is sent from the server to the client to return the
// definition of a symbol at a given text document position.
type DefinitionResponse struct {
	Response
	Result Location `json:"result"`
}
