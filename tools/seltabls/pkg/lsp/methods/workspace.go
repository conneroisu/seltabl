package methods

// Workspace Methods
//
// The methods in this package are used by the LSP to interact with the
// workspace.
// Interacting with the workspace is defined in the LSP specification.
//
// It includes methods for interacting with the workspace such as
// workspace/didChangeConfiguration, workspace/didChangeWatchedFiles, workspace/symbol, and workspace/executeCommand.
const (
	// MethodWorkspaceDidChangeConfiguration is the workspace did change
	// configuration notification method for the language server protocol.
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_didChangeConfiguration
	MethodWorkspaceDidChangeConfiguration Method = "workspace/didChangeConfiguration"

	// MethodWorkspaceDidChangeWatchedFiles is the workspace did change
	// watched files method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_didChangeWatchedFiles
	MethodWorkspaceDidChangeWatchedFiles Method = "workspace/didChangeWatchedFiles"

	// MethodWorkspaceSymbol is the workspace symbol method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_symbol
	MethodWorkspaceSymbol Method = "workspace/symbol"

	// MethodWorkspaceExecuteCommand is the workspace execute command method for the LSP
	//
	// Microsoft LSP Docs:
	// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_executeCommand
	MethodWorkspaceExecuteCommand Method = "workspace/executeCommand"
)
