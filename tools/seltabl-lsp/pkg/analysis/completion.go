package analysis

import (
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
)

var (
	// headerTag is the tag used to match a header cell's Value.
	headerTag = lsp.CompletionItem{
		Label:         "seltabl",
		Detail:        "Title Text for the header",
		Documentation: "This is the documentation for the header",
	}
	// selectorDataTag is the tag used to mark a data cell.
	selectorDataTag = lsp.CompletionItem{
		Label:         "dSel",
		Detail:        "Title Text for the data selector",
		Documentation: "This is the documentation for the data selector",
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorHeaderTag = lsp.CompletionItem{
		Label:         "hSel",
		Detail:        "Title Text for the header selector",
		Documentation: "This is the documentation for the header selector",
	}
	// selectorQueryTag is the tag used to signify selecting aspects of a cell
	selectorQueryTag = lsp.CompletionItem{
		Label:         "qSel",
		Detail:        "Title Text for the query selector",
		Documentation: "This is the documentation for the query selector",
	}
	// selectorMustBePresentTag is the tag used to signify selecting aspects of a cell
	selectorMustBePresentTag = lsp.CompletionItem{
		Label:         "must",
		Detail:        "Title Text for the must be present selector",
		Documentation: "This is the documentation for the must be present selector",
	}
	// selectorControlTag is the tag used to signify selecting aspects of a cell
	selectorControlTag = lsp.CompletionItem{
		Label:         "ctl",
		Detail:        "Title Text for the control selector",
		Documentation: "This is the documentation for the control selector",
	}
	// selectors is the slice of selectors to return for completions inside a struct tag but not a "" selector
	selectors = []lsp.CompletionItem{
		headerTag,
		selectorDataTag,
		selectorHeaderTag,
		selectorQueryTag,
		selectorMustBePresentTag,
		selectorControlTag,
	}
)

// TextDocumentCompletion returns the completions for a given text document.
func (s *State) TextDocumentCompletion(
	id int,
	document *lsp.TextDocumentIdentifier,
	_ *lsp.Position,
) lsp.CompletionResponse {
	text := s.Documents[document.URI]
	urls, _, err := s.getUrlsAndIgnores(text)
	if err != nil {
		s.Logger.Printf("failed to get urls and ignores: %s\n", err)
	}
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "Very cool editor",
			Documentation: "Fun to watch in videos. Don't forget to like & subscribe to streamers using it :)",
		},
	}
	for _, url := range urls {
		selectors := s.Selectors[url]
		for _, selector := range selectors {
			items = append(items, lsp.CompletionItem{
				Label:         selector.Selector,
				Detail:        "from: " + selector.URL.URL,
				Documentation: "A selector for the " + selector.Selector,
			})
		}
	}
	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
	return response
}
