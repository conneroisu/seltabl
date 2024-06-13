package lsp

// TextDocumentItem is a text document.
type TextDocumentItem struct {
	/**
	 * The text document's URI.
	 */
	URI string `json:"uri"`

	/**
	 * The text document's language identifier.
	 */
	LanguageID string `json:"languageId"`

	/**
	 * The version number of this document (it will increase after each
	 * change, including undo/redo).
	 */
	Version int `json:"version"`

	/**
	 * The content of the opened text document.
	 */
	Text string `json:"text"`
}

// TextDocumentIdentifier is a unique identifier for a text document.
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// VersionTextDocumentIdentifier is a text document with a version number.
type VersionTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// TextDocumentPositionParams is a text document position parameters.
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// Position is a position inside a text document.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Location is a location inside a resource, such as a line
// inside a text file.
type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// Range is a range in a text document.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// WorkspaceEdit is the workspace edit object.
type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes"`
}

// TextEdit represents an edit operation on a single text document.
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}
