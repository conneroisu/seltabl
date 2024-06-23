package lsp

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
	// Text is a completion item kind
	Text CompletionItemKind = iota + 1
	// Method is a completion item kind for a method or function completion
	Method
	// Function is a completion item kind for a function completion
	Function
	// Constructor is a completion item kind for a constructor completion
	Constructor
	// Field is a completion item kind for a field completion
	Field
	// Variable is a completion item kind for a variable completion
	Variable
	// Class is a completion item kind for a class completion
	Class
	// Interface is a completion item kind for an interface completion
	Interface
	// Module is a completion item kind for a module completion
	Module
	// Property is a completion item kind for a property completion
	Property
	// Unit is a completion item kind for a unit
	Unit
	// Value is a completion item kind for a value
	Value
	// Enum is a completion item kind for an enum
	Enum
	// Keyword is a completion item kind for a keyword
	Keyword
	// Snippet is a completion item kind for a snippet
	Snippet
	// Color is a completion item kind for a color
	Color
	// File is a completion item kind for a file
	File
	// Reference is a completion item kind for a reference
	Reference
	// Folder is a completion item kind for a folder
	Folder
	// EnumMember is a completion item kind for an enum member
	EnumMember
	// Constant is a completion item kind for a constant
	Constant
	// Struct is a completion item kind for a struct
	Struct
	// Event is a completion item kind for an event
	Event
	// Operator is a completion item kind for an operator
	Operator
	// TypeParameter is a completion item kind for a type parameter
	TypeParameter
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
