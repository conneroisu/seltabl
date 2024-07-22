package lsp

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// NotificationDidOpenTextDocument is a notification that is sent when
// the client opens a text document.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didOpen
type NotificationDidOpenTextDocument struct {
	// DidOpenTextDocumentNotification embeds the Notification struct
	Notification
	// Params are the parameters for the notification.
	Params protocol.DidOpenTextDocumentParams `json:"params"`
}

// Method returns the method for the did open text document params notification.
func (r NotificationDidOpenTextDocument) Method() methods.Method {
	return methods.MethodRequestTextDocumentDidOpen
}

// PublishDiagnosticsNotification is the notification for publishing diagnostics.
type PublishDiagnosticsNotification struct {
	// PublishDiagnosticsNotification embeeds the notification struct.
	Notification
	// Params are the parameters for the publish diagnostics notification.
	Params protocol.PublishDiagnosticsParams `json:"params"`
}

// Method returns the method for the publish diagnostics notification
func (r PublishDiagnosticsNotification) Method() methods.Method {
	return methods.NotificationPublishDiagnostics
}

const (
	// RPCVersion is the version of the RPC protocol.
	RPCVersion = "2.0"
)

// DidSaveTextDocumentNotification is a notification for when
// the client saves a text document.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didSave
type DidSaveTextDocumentNotification struct {
	// DidSaveTextDocumentParams embeds the Notification struct
	Notification
	// Params are the parameters for the notification.
	Params protocol.DidSaveTextDocumentParams `json:"params"`
}

// WillSaveTextDocumentNotification is a struct for the will save text document notification
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_willSave
type WillSaveTextDocumentNotification struct {
	Notification
	Params protocol.WillSaveTextDocumentParams `json:"params"`
}

// Method returns the method for the did save text document params notification
func (r DidSaveTextDocumentNotification) Method() methods.Method {
	return methods.MethodNotificationTextDocumentDidSave
}

// DidCloseTextDocumentParamsNotification is a struct for the did close text document params notification
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didClose
type DidCloseTextDocumentParamsNotification struct {
	Notification
	Params protocol.DidCloseTextDocumentParams `json:"params"`
}

// Method returns the method for the did close text document params notification
func (r DidCloseTextDocumentParamsNotification) Method() methods.Method {
	return methods.NotificationTextDocumentDidClose
}

// NewDidCloseTextDocumentParamsNotification returns a new did close text document params notification
func NewDidCloseTextDocumentParamsNotification(
	uri protocol.DocumentURI,
) DidCloseTextDocumentParamsNotification {
	return DidCloseTextDocumentParamsNotification{
		Notification: Notification{
			RPC: RPCVersion,
		},
		Params: protocol.DidCloseTextDocumentParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		},
	}
}

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
