package lsp

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// CompletionRequest is a request for a completion to the language server
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type CompletionRequest struct {
	// CompletionRequest embeds the Request struct
	Request
	// Params are the parameters for the completion request
	Params protocol.CompletionParams `json:"params"`
}

// Method returns the method for the completion request
func (r CompletionRequest) Method() methods.Method {
	return methods.MethodRequestTextDocumentCompletion
}

// CompletionResponse is a response for a completion to the language server
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type CompletionResponse struct {
	// CompletionResponse embeds the Response struct
	Response
	// Result is the result of the completion request
	Result []protocol.CompletionItem `json:"result"`
}

// Method returns the method for the completion response
func (r CompletionResponse) Method() methods.Method {
	return methods.MethodRequestTextDocumentCompletion
}
