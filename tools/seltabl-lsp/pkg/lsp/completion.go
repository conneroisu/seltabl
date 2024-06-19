package lsp

// CompletionRequest is a request for a completion to the language server
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type CompletionRequest struct {
	// CompletionRequest embeds the Request struct
	Request
	// Params are the parameters for the completion request
	Params CompletionParams `json:"params"`
}

// CompletionParams is a struct for the completion params
type CompletionParams struct {
	// CompletionParams embeds the TextDocumentPositionParams struct
	TextDocumentPositionParams
}

// CompletionResponse is a response for a completion to the language server
type CompletionResponse struct {
	// CompletionResponse embeds the Response struct
	Response
	// Result is the result of the completion request
	Result []CompletionItem `json:"result"`
}

// CompletionItem is a struct for a completion item
type CompletionItem struct {
	// Label is the label for the completion item
	Label string `json:"label"`
	// Detail is the detail for the completion item
	Detail string `json:"detail"`
	// Documentation is the documentation for the completion item
	Documentation string `json:"documentation"`
}
