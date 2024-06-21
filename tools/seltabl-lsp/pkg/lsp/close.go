package lsp

// DidCloseTextDocumentParamsNotification is a struct for the did close text document params notification
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didClose
type DidCloseTextDocumentParamsNotification struct {
	Notification
	Params DidCloseTextDocumentParamsNotificationParams `json:"params"`
}

// NewDidCloseTextDocumentParamsNotification returns a new did close text document params notification
func NewDidCloseTextDocumentParamsNotification() DidCloseTextDocumentParamsNotification {
	return DidCloseTextDocumentParamsNotification{
		Notification: Notification{
			RPC: "2.0",
		},
	}
}

// DidCloseTextDocumentParamsNotificationParams is a struct for the did close text document params notification params
type DidCloseTextDocumentParamsNotificationParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
