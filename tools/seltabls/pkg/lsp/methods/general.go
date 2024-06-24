package methods

// Notification Methods
const (
	// MethodInitialize is the initialize notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
	MethodInitialize Method = "initialize"

	// MethodInitialized is the initialized notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialized
	MethodInitialized Method = "initialized"

	// MethodShutdown is the shutdown notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#shutdown
	MethodShutdown Method = "shutdown"

	// MethodExit is the exit notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#exit
	MethodExit Method = "exit"

	// MethodCancelRequest is the cancel request method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#cancelRequest
	MethodCancelRequest Method = "$/cancelRequest"
)
