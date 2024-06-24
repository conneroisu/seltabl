package methods

const (
	// MethodTextDocumentDidSave is the text document did save notification
	// method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didSave
	MethodTextDocumentDidSave Method = "textDocument/didSave"
)
