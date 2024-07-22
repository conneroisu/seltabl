package lsp

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// TextDocumentCompletionRequest is a request for a completion to the language server
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type TextDocumentCompletionRequest struct {
	// CompletionRequest embeds the Request struct
	Request
	// Params are the parameters for the completion request
	Params protocol.CompletionParams `json:"params"`
}

// Method returns the method for the completion request
func (r TextDocumentCompletionRequest) Method() methods.Method {
	return methods.MethodRequestTextDocumentCompletion
}

// TextDocumentCompletionResponse is a response for a completion to the language server
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type TextDocumentCompletionResponse struct {
	// CompletionResponse embeds the Response struct
	Response
	// Result is the result of the completion request
	Result []protocol.CompletionItem `json:"result"`
}

// Method returns the method for the completion response
func (r TextDocumentCompletionResponse) Method() methods.Method {
	return methods.MethodRequestTextDocumentCompletion
}
