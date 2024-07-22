package analysis

import (
	"go.lsp.dev/protocol"
)

var (
	// selectorDataTag is the tag used to mark a data cell.
	selectorDataTag = protocol.CompletionItem{Label: "dSel",
		Detail:        "Title Text for the data selector",
		Documentation: "The data selector is used to find the data column for the given field.",
		Kind:          protocol.CompletionItemKindField,
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorHeaderTag = protocol.CompletionItem{Label: "hSel",
		Detail:        "Title Text for the header selector",
		Documentation: "The header selector is used to find the header row and column for the given field. This selection is removed from the data selection.",
		Kind:          protocol.CompletionItemKindField,
	}
	// selectorQueryTag is the tag used to signify selecting aspects of a cell
	selectorQueryTag = protocol.CompletionItem{Label: "qSel",
		Detail:        "Title Text for the query selector",
		Documentation: "The query selector is used to find the data column for the given field. This selection is removed from the header selection.",
		Kind:          protocol.CompletionItemKindField,
	}
	// selectorMustBePresentTag is the tag used to signify selecting aspects of a cell
	selectorMustBePresentTag = protocol.CompletionItem{Label: "must",
		Detail:        "Title Text for the must be present selector",
		Documentation: "The must be present selector is used to ensure that the given raw text is present in the html content.",
		Kind:          protocol.CompletionItemKindField,
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorControlTag = protocol.CompletionItem{Label: "ctl",
		Detail:        "Title Text for the control selector",
		Documentation: "The control selector is used to control how the given document is parsed to get the value of the cell/field.",
		Kind:          protocol.CompletionItemKindField,
	}
	// completionKeys is the slice of completionKeys to return for completions inside a struct tag but not a "" selector
	completionKeys = []protocol.CompletionItem{
		selectorDataTag,
		selectorHeaderTag,
		selectorQueryTag,
		selectorMustBePresentTag,
		selectorControlTag,
	}
)
