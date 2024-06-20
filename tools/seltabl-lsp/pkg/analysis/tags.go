package analysis

import "github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"

var (
	// headerTag is the tag used to match a header cell's Value.
	headerTag = lsp.CompletionItem{Label: "seltabl",
		Detail:        "Title Text for the header",
		Documentation: "This is the documentation for the header",
	}
	// selectorDataTag is the tag used to mark a data cell.
	selectorDataTag = lsp.CompletionItem{Label: "dSel",
		Detail:        "Title Text for the data selector",
		Documentation: "This is the documentation for the data selector",
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorHeaderTag = lsp.CompletionItem{Label: "hSel",
		Detail:        "Title Text for the header selector",
		Documentation: "This is the documentation for the header selector",
	}
	// selectorQueryTag is the tag used to signify selecting aspects of a cell
	selectorQueryTag = lsp.CompletionItem{Label: "qSel",
		Detail:        "Title Text for the query selector",
		Documentation: "This is the documentation for the query selector",
	}
	// selectorMustBePresentTag is the tag used to signify selecting aspects of a cell
	selectorMustBePresentTag = lsp.CompletionItem{Label: "must",
		Detail:        "Title Text for the must be present selector",
		Documentation: "This is the documentation for the must be present selector",
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorControlTag = lsp.CompletionItem{Label: "ctl",
		Detail:        "Title Text for the control selector",
		Documentation: "This is the documentation for the control selector",
	}
)
