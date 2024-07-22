package methods

// Text Document Request Methods
const (
	// MethodRequestTextDocumentDidOpen is the text document did open
	// request method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didOpen
	MethodRequestTextDocumentDidOpen Method = "textDocument/didOpen"

	// MethodRequestTextDocumentCompletion is the text document completion
	// request method.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
	MethodRequestTextDocumentCompletion Method = "textDocument/completion"

	// MethodRequestTextDocumentHover is the text document hover request
	// method.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
	MethodRequestTextDocumentHover Method = "textDocument/hover"

	// MethodRequestTextDocumentSignatureHelp is the text document signature
	// help method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_signatureHelp
	MethodRequestTextDocumentSignatureHelp Method = "textDocument/signatureHelp"

	// MethodRequestTextDocumentDefinition is the text document definition
	// method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_definition
	MethodRequestTextDocumentDefinition Method = "textDocument/definition"

	// MethodTextDocumentReferences is the text document references method
	// for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_references
	MethodTextDocumentReferences Method = "textDocument/references"

	// MethodRequestTextDocumentDocumentHighlight is the text document
	// document highlight method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_documentHighlight
	MethodRequestTextDocumentDocumentHighlight Method = "textDocument/documentHighlight"

	// MethodRequestTextDocumentDocumentSymbol is the text document symbol
	// method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_documentSymbol
	MethodRequestTextDocumentDocumentSymbol Method = "textDocument/documentSymbol"

	// MethodTextDocumentFormatting is the text document formatting method
	// for the LSP.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_formatting
	MethodTextDocumentFormatting Method = "textDocument/formatting"

	// NotificationTextDocumentDidClose is the text document did close
	// notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didClose
	NotificationTextDocumentDidClose Method = "textDocument/didClose"

	// MethodTextDocumentRangeFormatting is the text document range
	// formatting method for the LSP.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_rangeFormatting
	MethodTextDocumentRangeFormatting Method = "textDocument/rangeFormatting"

	// MethodTextDocumentOnTypeFormatting is the text document on type
	// formatting method for the LSP.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_onTypeFormatting
	MethodTextDocumentOnTypeFormatting Method = "textDocument/onTypeFormatting"

	// MethodTextDocumentRename is the text document code action method for
	// the LSP.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_rename
	MethodTextDocumentRename Method = "textDocument/rename"

	// MethodRequestTextDocumentCodeAction is the text document code action
	// method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_codeAction
	MethodRequestTextDocumentCodeAction Method = "textDocument/codeAction"

	// MethodTextDocumentCodeLens is the text document code lens method for
	// the LSP.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_codeLens
	MethodTextDocumentCodeLens Method = "textDocument/codeLens"

	// MethodTextDocumentDocumentLink is the text document document link
	// method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_documentLink
	MethodTextDocumentDocumentLink Method = "textDocument/documentLink"
)

// Notification methods.
const (
	// NotificationPublishDiagnostics is the publish diagnostics
	// notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_publishDiagnostics
	NotificationPublishDiagnostics Method = "textDocument/publishDiagnostics"

	// MethodNotificationTextDocumentDidSave is the text document did save
	// notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didSave
	MethodNotificationTextDocumentDidSave Method = "textDocument/didSave"

	// MethodNotificationTextDocumentWillSave is the text document will save
	// notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_willSave
	MethodNotificationTextDocumentWillSave Method = "textDocument/willSave"

	// NotificationDidSaveTextDocument is the text document did save
	// notification for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didSave
	NotificationDidSaveTextDocument = "textDocument/didSave"

	// NotificationMethodTextDocumentDidChange is the text document did
	// change request method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
	NotificationMethodTextDocumentDidChange Method = "textDocument/didChange"

	// NotificationMethodLogMessage is the log message notification method
	// for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_logMessage
	NotificationMethodLogMessage Method = "window/logMessage"
)
