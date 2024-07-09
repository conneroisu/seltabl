package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// DiagnosticSeverity is an enum for diagnostic severities.
type DiagnosticSeverity int

const (
	// DiagnosticError reports an error.
	DiagnosticError DiagnosticSeverity = iota + 1
	// DiagnosticWarning reports a warning.
	DiagnosticWarning
	// DiagnosticInformation reports an information.
	DiagnosticInformation
	// DiagnosticHint reports a hint.
	DiagnosticHint
)

// PublishDiagnosticsNotification is the notification for publishing diagnostics.
type PublishDiagnosticsNotification struct {
	// PublishDiagnosticsNotification embeeds the notification struct.
	Notification
	// Params are the parameters for the publish diagnostics notification.
	Params PublishDiagnosticsParams `json:"params"`
}

// Method returns the method for the publish diagnostics notification
func (r PublishDiagnosticsNotification) Method() methods.Method {
	return methods.MethodNotificationTextDocumentDidSave
}

// PublishDiagnosticsParams are the parameters for the publish diagnostics notification.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_publishDiagnostics
type PublishDiagnosticsParams struct {
	// URI is the uri for the diagnostics.
	URI string `json:"uri"`
	// Diagnostics are the diagnostics for the uri.
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// Diagnostic is a struct for a diagnostic.
type Diagnostic struct {
	// Range is the range for the diagnostic.
	Range Range `json:"range"`
	// Severity is the severity for the diagnostic.
	Severity DiagnosticSeverity `json:"severity"`
	// Source is the source for the diagnostic.
	Source string `json:"source"`
	// Message is the message for the diagnostic.
	Message string `json:"message"`
}

// String returns the string representation of the DiagnosticSeverity.
func (d DiagnosticSeverity) String() string {
	return [...]string{
		"Error",
		"Warning",
		"Information",
		"Hint",
	}[d-1]
}
