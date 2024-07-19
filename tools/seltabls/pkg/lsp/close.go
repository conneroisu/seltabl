package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// DidCloseTextDocumentParamsNotification is a struct for the did close text document params notification
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didClose
type DidCloseTextDocumentParamsNotification struct {
	Notification
	Params DidCloseTextDocumentParamsNotificationParams `json:"params"`
}

// Method returns the method for the did close text document params notification
func (r DidCloseTextDocumentParamsNotification) Method() methods.Method {
	return methods.NotificationTextDocumentDidClose
}

// NewDidCloseTextDocumentParamsNotification returns a new did close text document params notification
func NewDidCloseTextDocumentParamsNotification(
	uri string,
) DidCloseTextDocumentParamsNotification {
	return DidCloseTextDocumentParamsNotification{
		Notification: Notification{
			RPC: RPCVersion,
		},
		Params: DidCloseTextDocumentParamsNotificationParams{
			TextDocument: TextDocumentIdentifier{URI: uri},
		},
	}
}

// DidCloseTextDocumentParamsNotificationParams is a struct for the did close text document params notification params
type DidCloseTextDocumentParamsNotificationParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
