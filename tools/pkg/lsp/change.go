package lsp

// TextDocumentDidChangeNotification is sent from the client to the server to signal
// that the content of a text document has changed.
type TextDocumentDidChangeNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

// DidChangeTextDocumentParams is sent from the client to the server to signal
// that the content of a text document has changed.
type DidChangeTextDocumentParams struct {
	// TextDocument is the  document that did change. The version number points to the version
	TextDocument VersionTextDocumentIdentifier `json:"textDocument"`
	// ContentChanges is the array of content changes.
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// TextDocumentContentChangeEvent is sent from the client to the server to signal
// that the content of a text document has changed.
type TextDocumentContentChangeEvent struct {
	// The new text of the whole document.
	Text string `json:"text"`
}
