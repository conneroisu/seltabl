package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// CodeActionRequest is a request for a code action to the language server.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_code
type CodeActionRequest struct {
	// CodeActionRequest embeds the Request struct
	Request
	// Params are the parameters for the code action request.
	Params TextDocumentCodeActionParams `json:"params"`
}

// Method returns the method for the code action request
func (r CodeActionRequest) Method() methods.Method {
	return methods.MethodRequestTextDocumentCodeAction
}

// TextDocumentCodeActionParams are the parameters for a code action request.
type TextDocumentCodeActionParams struct {
	// TextDocument is the text document for the code action request.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// Range is the range for the code action request.
	Range Range `json:"range"`
	// Context is the context for the code action request.
	Context CodeActionContext `json:"context"`
}

// TextDocumentCodeActionResponse is the response for a code action request.
type TextDocumentCodeActionResponse struct {
	// TextDocumentCodeActionResponse embeds the Response struct
	Response
	// Result is the result for the code action request.
	Result []CodeAction `json:"result"`
}

// Method returns the method for the code action response
func (r TextDocumentCodeActionResponse) Method() methods.Method {
	return methods.MethodRequestTextDocumentCodeAction
}

// CodeActionContext is the context for a code action request.
type CodeActionContext struct {
	// Add fields for CodeActionContext as needed
}

// CodeAction is a code action for a given text document.
type CodeAction struct {
	// Title is the title for the code action.
	Title string `json:"title"`
	// Edit is the edit for the code action.
	Edit *WorkspaceEdit `json:"edit,omitempty"`
	// Command is the command for the code action.
	Command *Command `json:"command,omitempty"`
}

// Command is a command for a given text document.
type Command struct {
	// Title is the title for the command.
	Title string `json:"title"`
	// Command is the command for the command.
	Command string `json:"command"`
	// Arguments are the arguments for the command.
	Arguments []interface{} `json:"arguments,omitempty"`
}
