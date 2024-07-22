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
