package lsp

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// TextDocumentDidChangeNotification is sent from the client to the server to signal
// that the content of a text document has changed.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
type TextDocumentDidChangeNotification struct {
	// TextDocumentDidChangeNotification embeds the Notification struct
	Notification
	// Params are the parameters for the notification.
	Params protocol.DidChangeTextDocumentParams `json:"params"`
}

// Method returns the method for the text document did change notification
func (r TextDocumentDidChangeNotification) Method() methods.Method {
	return methods.NotificationMethodTextDocumentDidChange
}
