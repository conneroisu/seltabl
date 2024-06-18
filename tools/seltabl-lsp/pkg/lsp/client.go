package lsp

// ClientCapabilities is a struct for the client capabilities
type ClientCapabilities struct {
	// TextDocumentSync is the text document sync for the client capabilities.
	TextDocumentSync int `json:"textDocumentSync"`
	// HoverProvider is a boolean indicating whether the server provides.
	HoverProvider bool `json:"hoverProvider"`
	// DefinitionProvider is a boolean indicating whether the server provides definition capabilities.
	DefinitionProvider bool `json:"definitionProvider"`
	// CodeActionProvider is a boolean indicating whether the server provides code actions.
	CodeActionProvider bool `json:"codeActionProvider"`
	// CompletionProvider is a map of completion providers.
	CompletionProvider map[string]any `json:"completionProvider"`
}
