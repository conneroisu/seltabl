package methods

// General Notification Methods
const (
	// MethodInitialize is the initialize notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
	MethodInitialize Method = "initialize"

	// MethodNotificationInitialized is the initialized notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialized
	MethodNotificationInitialized Method = "initialized"

	// MethodNotificationExit is the exit notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#exit
	MethodNotificationExit Method = "exit"
)

// General Request Methods
const (
	// MethodCancelRequest is the cancel request method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#cancelRequest
	MethodCancelRequest Method = "$/cancelRequest"

	// MethodShutdown is the shutdown notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
	MethodShutdown Method = "shutdown"
)
