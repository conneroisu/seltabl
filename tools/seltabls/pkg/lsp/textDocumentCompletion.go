package lsp

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"

// CompletionRequest is a request for a completion to the language server
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type CompletionRequest struct {
	// CompletionRequest embeds the Request struct
	Request
	// Params are the parameters for the completion request
	Params CompletionParams `json:"params"`
}

// Method returns the method for the completion request
func (r CompletionRequest) Method() methods.Method {
	return methods.MethodRequestTextDocumentCompletion
}

// CompletionParams is a struct for the completion params
type CompletionParams struct {
	// CompletionParams embeds the TextDocumentPositionParams struct
	TextDocumentPositionParams
}

// CompletionResponse is a response for a completion to the language server
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type CompletionResponse struct {
	// CompletionResponse embeds the Response struct
	Response
	// Result is the result of the completion request
	Result []CompletionItem `json:"result"`
}

// Method returns the method for the completion response
func (r CompletionResponse) Method() string {
	return "textDocument/completion"
}

// CompletionItem is a struct for a completion item
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_completion
type CompletionItem struct {
	// Label is the label for the completion item
	Label string `json:"label"`
	// Detail is the detail for the completion item
	Detail string `json:"detail"`
	// Documentation is the documentation for the completion item
	Documentation string `json:"documentation"`
	// Kind is the kind of the completion item
	Kind CompletionItemKind `json:"kind"`
}

// CompletionItemTag is a struct for a completion item tag
type CompletionItemTag struct {
	Deprecated bool `json:"deprecated"`
}

// CompletionItemKind is an enum for completion item kinds.
type CompletionItemKind int

const (
	// CompletionKindText is a completion item kind
	CompletionKindText CompletionItemKind = iota + 1
	// CompletionKindMethod is a completion item kind for a method or function completion
	CompletionKindMethod
	// CompletionKindFunction is a completion item kind for a function completion
	CompletionKindFunction
	// CompletionKindConstructor is a completion item kind for a constructor completion
	CompletionKindConstructor
	// CompletionKindField is a completion item kind for a field completion
	CompletionKindField
	// CompletionKindVariable is a completion item kind for a variable completion
	CompletionKindVariable
	// CompletionKindClass is a completion item kind for a class completion
	CompletionKindClass
	// CompletionKindInterface is a completion item kind for an interface completion
	CompletionKindInterface
	// CompletionKindModule is a completion item kind for a module completion
	CompletionKindModule
	// CompletionKindProperty is a completion item kind for a property completion
	CompletionKindProperty
	// CompletionKindUnit is a completion item kind for a unit
	CompletionKindUnit
	// CompletionKindValue is a completion item kind for a value
	CompletionKindValue
	// CompletionKindEnum is a completion item kind for an enum
	CompletionKindEnum
	// CompletionKindKeyword is a completion item kind for a keyword
	CompletionKindKeyword
	// CompletionKindSnippet is a completion item kind for a snippet
	CompletionKindSnippet
	// CompletionKindColor is a completion item kind for a color
	CompletionKindColor
	// CompletionKindFile is a completion item kind for a file
	CompletionKindFile
	// CompletionKindReference is a completion item kind for a reference
	CompletionKindReference
	// CompletionKindFolder is a completion item kind for a folder
	CompletionKindFolder
	// CompletionKindEnumMember is a completion item kind for an enum member
	CompletionKindEnumMember
	// CompletionKindConstant is a completion item kind for a constant
	CompletionKindConstant
	// CompletionKindStruct is a completion item kind for a struct
	CompletionKindStruct
	// CompletionindEvent is a completion item kind for an event
	CompletionindEvent
	// CompletionKindOperator is a completion item kind for an operator
	CompletionKindOperator
	// CompletionKindTypeParameter is a completion item kind for a type parameter
	CompletionKindTypeParameter
)

// String returns the string representation of the CompletionItemKind.
func (c CompletionItemKind) String() string {
	return [...]string{
		"Text",
		"Method",
		"Function",
		"Constructor",
		"Field",
		"Variable",
		"Class",
		"Interface",
		"Module",
		"Property",
		"Unit",
		"Value",
		"Enum",
		"Keyword",
		"Snippet",
		"Color",
		"File",
		"Reference",
		"Folder",
		"EnumMember",
		"Constant",
		"Struct",
		"Event",
		"Operator",
		"TypeParameter",
	}[c-1]
}
