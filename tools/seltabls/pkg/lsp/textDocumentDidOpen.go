package lsp

// NotificationDidOpenTextDocument is a notification that is sent when
// the client opens a text document.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didOpen
type NotificationDidOpenTextDocument struct {
	// DidOpenTextDocumentNotification embeds the Notification struct
	Notification
	// Params are the parameters for the notification.
	Params DidOpenTextDocumentParams `json:"params"`
}

// DidOpenTextDocumentParams contains the text document after it has been opened.
type DidOpenTextDocumentParams struct {
	// TextDocument is the text document after it has been opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}
