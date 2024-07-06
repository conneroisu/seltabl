package methods

const (
	// MethodNotificationTextDocumentDidSave is the text document did save notification
	// method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didSave
	MethodNotificationTextDocumentDidSave Method = "textDocument/didSave"
)
