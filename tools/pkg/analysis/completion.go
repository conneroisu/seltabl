package analysis

import (
	"strconv"

	"github.com/conneroisu/seltabl/tools/data"
	"github.com/conneroisu/seltabl/tools/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/pkg/parsers"
)

var (
	headerTag = lsp.CompletionItem{
		Label:         "seltabl",
		Detail:        "Title Text for the header",
		Documentation: "This is the documentation for the header",
	}
	// headerTag is the tag used to match a header cell's Value.
	// selectorDataTag is the tag used to mark a data cell.
	selectorDataTag = lsp.CompletionItem{
		Label:         "dSel",
		Detail:        "Title Text for the data selector",
		Documentation: "This is the documentation for the data selector",
	}
	selectorHeaderTag = lsp.CompletionItem{
		Label:         "hSel",
		Detail:        "Title Text for the header selector",
		Documentation: "This is the documentation for the header selector",
	}
	selectorQueryTag = lsp.CompletionItem{
		Label:         "qSel",
		Detail:        "Title Text for the query selector",
		Documentation: "This is the documentation for the query selector",
	}
	selectorMustBePresentTag = lsp.CompletionItem{
		Label:         "must",
		Detail:        "Title Text for the must be present selector",
		Documentation: "This is the documentation for the must be present selector",
	}
	selectorControlTag = lsp.CompletionItem{
		Label:         "ctl",
		Detail:        "Title Text for the control selector",
		Documentation: "This is the documentation for the control selector",
	}
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
	location *lsp.Position,
) lsp.CompletionResponse {
	src := s.Documents[document.URI]
	s.Logger.Println("Received text document completion uri: " + document.URI)
	pos, err := parsers.PositionStatusInStructTag(src, *location)
	s.Logger.Println("Position Status in struct tag: " + strconv.Itoa(pos) + " for location: line: " + strconv.Itoa(location.Line) + "character: " + strconv.Itoa(location.Character))
	if err != nil {
		s.Logger.Printf("failed to get position status in struct tag: %s\n", err)
		s.Logger.Printf("src: %s\n", src)
		return lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: nil,
		}
	}
	if pos == 0 {
		s.Logger.Println("OUTSIDE: POSITION STATUS IN STRUCT TAG")
		return lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: nil,
		}
	}
	if pos == 1 {
		s.Logger.Println("INSIDE: POSITION STATUS IN STRUCT TAG")
		return lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: selectors,
		}
	}
	s.Logger.Println("OUTSIDE: POSITION STATUS INSIDE STRUCT TAG AND \"\"")
	urls, err := parsers.ExtractUrls(src)
	if err != nil {
		s.Logger.Printf("failed to extract urls: %s\n", err)
		return lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: nil,
		}
	}
	ignores, err := parsers.ExtractIgnores(src)
	var selectors []data.Selector
	for _, url := range urls {
		got, err := s.getSelectors(url, ignores)
		if err != nil {
			s.Logger.Printf("failed to get selectors: %s\n", err)
		}
		selectors = append(selectors, got...)
	}
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "Very cool editor",
			Documentation: "Fun to watch in videos. Don't forget to like & subscribe to streamers using it :)",
		},
	}
	for _, selector := range selectors {
		items = append(items, lsp.CompletionItem{
			Label:         selector.Selector,
			Detail:        "from: " + selector.URL,
			Documentation: "A selector for the " + selector.Selector,
		})
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
