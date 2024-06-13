package lsp

// DidOpenTextDocumentNotification is a notification that is sent when
// the client opens a text document.
type DidOpenTextDocumentNotification struct {
	Notification
	Params DidOpenTextDocumentParams `json:"params"`
}

// DidOpenTextDocumentParams contains the text document after it has been opened.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}
