package lsp

// DidOpenTextDocumentNotification is a notification that is sent when
// the client opens a text document.
type DidOpenTextDocumentNotification struct {
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
