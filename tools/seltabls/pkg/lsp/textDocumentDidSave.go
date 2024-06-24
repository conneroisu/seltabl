package lsp

// DidSaveTextDocumentParamsNotification is a notification for when
// the client saves a text document.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didSave
type DidSaveTextDocumentParamsNotification struct {
	// DidSaveTextDocumentParams embeds the Notification struct
	Notification
	// Params are the parameters for the notification.
	Params DidSaveTextDocumentParams `json:"params"`
}

// DidSaveTextDocumentParams contains the text document after it has been saved.
type DidSaveTextDocumentParams struct {
	// TextDocument is the text document after it has been saved.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
