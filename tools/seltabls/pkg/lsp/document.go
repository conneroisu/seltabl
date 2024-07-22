package lsp

// TextDocumentIdentifier is a unique identifier for a text document.
type TextDocumentIdentifier struct {
	// URI is the uri for the text document.
	URI string `json:"uri"`
}
