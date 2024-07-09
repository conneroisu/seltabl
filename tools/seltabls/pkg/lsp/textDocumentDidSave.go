package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

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

// Method returns the method for the did save text document params notification
func (r DidSaveTextDocumentParamsNotification) Method() methods.Method {
	return methods.MethodNotificationTextDocumentDidSave
}

// DidSaveTextDocumentParams contains the text document after it has been saved.
type DidSaveTextDocumentParams struct {
	// TextDocument is the text document after it has been saved.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
