package methods

// Text Document Notification Methods
const (
	// NotificationDidSaveTextDocument is the text document did save notification for the LSP
	NotificationDidSaveTextDocument = "textDocument/didSave"

	// NotificationMethodTextDocumentDidChange is the text document did change request method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
	NotificationMethodTextDocumentDidChange Method = "textDocument/didChange"

	// NotificationMethodLogMessage is the log message notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_logMessage
	NotificationMethodLogMessage Method = "window/logMessage"
)
