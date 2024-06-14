package lsp

// DefinitionRequest is sent from the client to the server to resolve a
// definition.
type DefinitionRequest struct {
	// DefinitionRequest embeds the Request struct
	Request
	// Params are the parameters for the definition request
	Params DefinitionParams `json:"params"`
}

// DefinitionParams is the set of parameters that are sent in a
// DefinitionRequest.
type DefinitionParams struct {
	// DefinitionParams embeds the TextDocumentPositionParams struct
	TextDocumentPositionParams
}

// DefinitionResponse is sent from the server to the client to return the
// definition of a symbol at a given text document position.
type DefinitionResponse struct {
	// DefinitionResponse embeds the Response struct
	Response
	// Result is the result of the definition request
	Result Location `json:"result"`
}
