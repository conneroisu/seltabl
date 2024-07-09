package lsp

import "fmt"

// TextDocumentItem is a text document.
type TextDocumentItem struct {
	// URI is the uri for the text document.
	URI string `json:"uri"`

	// LanguageID is the language id for the text document.
	LanguageID string `json:"languageId"`

	// Version is the version number of a given text document.
	Version int `json:"version"`

	// Text is the text of the text document.
	Text string `json:"text"`
}

// TextDocumentIdentifier is a unique identifier for a text document.
type TextDocumentIdentifier struct {
	// URI is the uri for the text document.
	URI string `json:"uri"`
}

// VersionTextDocumentIdentifier is a text document with a version number.
type VersionTextDocumentIdentifier struct {
	// VersionTextDocumentIdentifier embeds the TextDocumentIdentifier struct
	TextDocumentIdentifier
	// Version is the version number for the text document.
	Version int `json:"version"`
}

// TextDocumentPositionParams is a text document position parameters.
type TextDocumentPositionParams struct {
	// TextDocument is the text document for the position parameters.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// Position is the position for the text document.
	Position Position `json:"position"`
}

// Position is a position inside a text document.
type Position struct {
	// Line is the line number for the position (zero-based).
	Line int `json:"line"`
	// Character is the character number for the position (zero-based).
	Character int `json:"character"`
}

// String returns a string representation of the position.
func (p Position) String() string {
	return fmt.Sprintf("Line: %d, Character: %d", p.Line, p.Character)
}

// Location is a location inside a resource, such as a line
// inside a text file.
type Location struct {
	// URI is the uri for the location.
	URI string `json:"uri"`
	// Range is the range for the location.
	Range Range `json:"range"`
}

// Range is a range in a text document.
type Range struct {
	// Start is the start of a given range.
	Start Position `json:"start"`
	// End is the end of a given range.
	End Position `json:"end"`
}

// WorkspaceEdit is the workspace edit object.
type WorkspaceEdit struct {
	// Changes is the changes for the workspace edit.
	Changes map[string][]TextEdit `json:"changes"`
}

// TextEdit represents an edit operation on a single text document.
type TextEdit struct {
	// Range is the range for the text edit.
	Range Range `json:"range"`
	// NewText is the new text for the text edit.
	NewText string `json:"newText"`
}

// LineRange returns a range of a line in a document
//
// line is the line number
//
// start is the start character of the range
//
// end is the end character of the range
func LineRange(line, start, end int) Range {
	return Range{
		Start: Position{
			Line:      line,
			Character: start,
		},
		End: Position{
			Line:      line,
			Character: end,
		},
	}
}
