package lsp

// DidCloseTextDocumentParamsNotification is a struct for the did close text document params notification
type DidCloseTextDocumentParamsNotification struct {
	Notification
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
