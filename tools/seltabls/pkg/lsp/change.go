package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// TextDocumentDidChangeNotification is sent from the client to the server to signal
// that the content of a text document has changed.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
type TextDocumentDidChangeNotification struct {
	// TextDocumentDidChangeNotification embeds the Notification struct
	Notification
	// Params are the parameters for the notification.
	Params DidChangeTextDocumentParams `json:"params"`
}

// Method returns the method for the text document did change notification
func (r TextDocumentDidChangeNotification) Method() methods.Method {
	return methods.NotificationMethodTextDocumentDidChange
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
