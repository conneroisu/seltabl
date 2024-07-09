package analysis

import "github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"

var (
	// selectorDataTag is the tag used to mark a data cell.
	selectorDataTag = lsp.CompletionItem{Label: "dSel",
		Detail:        "Title Text for the data selector",
		Documentation: "The data selector is used to find the data column for the given field.",
		Kind:          lsp.CompletionKindField,
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorHeaderTag = lsp.CompletionItem{Label: "hSel",
		Detail:        "Title Text for the header selector",
		Documentation: "The header selector is used to find the header row and column for the given field. This selection is removed from the data selection.",
		Kind:          lsp.CompletionKindField,
	}
	// selectorQueryTag is the tag used to signify selecting aspects of a cell
	selectorQueryTag = lsp.CompletionItem{Label: "qSel",
		Detail:        "Title Text for the query selector",
		Documentation: "The query selector is used to find the data column for the given field. This selection is removed from the header selection.",
		Kind:          lsp.CompletionKindField,
	}
	// selectorMustBePresentTag is the tag used to signify selecting aspects of a cell
	selectorMustBePresentTag = lsp.CompletionItem{Label: "must",
		Detail:        "Title Text for the must be present selector",
		Documentation: "The must be present selector is used to ensure that the given raw text is present in the html content.",
		Kind:          lsp.CompletionKindField,
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorControlTag = lsp.CompletionItem{Label: "ctl",
		Detail:        "Title Text for the control selector",
		Documentation: "The control selector is used to control how the given document is parsed to get the value of the cell/field.",
		Kind:          lsp.CompletionKindField,
	}
	// completionKeys is the slice of completionKeys to return for completions inside a struct tag but not a "" selector
	completionKeys = []lsp.CompletionItem{
		selectorDataTag,
		selectorHeaderTag,
		selectorQueryTag,
		selectorMustBePresentTag,
		selectorControlTag,
	}
)
