package lsp

// PublishDiagnosticsNotification is the notification for publishing diagnostics.
type PublishDiagnosticsNotification struct {
	// PublishDiagnosticsNotification embeeds the notification struct.
	Notification
	// Params are the parameters for the publish diagnostics notification.
	Params PublishDiagnosticsParams `json:"params"`
}

// PublishDiagnosticsParams are the parameters for the publish diagnostics notification.
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
	Severity int `json:"severity"`
	// Source is the source for the diagnostic.
	Source string `json:"source"`
	// Message is the message for the diagnostic.
	Message string `json:"message"`
}
