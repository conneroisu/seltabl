package lsp

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// TextDocumentCodeActionRequest is a request for a code action to the language server.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_codeAction
type TextDocumentCodeActionRequest struct {
	// CodeActionRequest embeds the Request struct
	Request
	// Params are the parameters for the code action request.
	Params protocol.CodeActionParams `json:"params"`
}

// Method returns the method for the code action request
func (r TextDocumentCodeActionRequest) Method() methods.Method {
	return methods.MethodRequestTextDocumentCodeAction
}

// TextDocumentCodeActionResponse is the response for a code action request.
type TextDocumentCodeActionResponse struct {
	// TextDocumentCodeActionResponse embeds the Response struct
	Response
	// Result is the result for the code action request.
	Result []protocol.CodeAction `json:"result"`
}

// Method returns the method for the code action response
func (r TextDocumentCodeActionResponse) Method() methods.Method {
	return methods.MethodRequestTextDocumentCodeAction
}
