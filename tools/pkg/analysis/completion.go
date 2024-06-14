package analysis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/conneroisu/seltabl/tools/data/master"
	"github.com/conneroisu/seltabl/tools/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/pkg/minify"
	"github.com/conneroisu/seltabl/tools/pkg/parsers"
)

// getLogger returns a logger that writes to a file
func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open file: %w", err))
	}
	return log.New(logFile, "[seltabl-lsp]", log.LstdFlags)
}

// TextDocumentCompletion returns the completions for a given text document.
func (s *State) TextDocumentCompletion(
	id int,
	document *lsp.TextDocumentIdentifier,
	location *lsp.Position,
) lsp.CompletionResponse {
	s.Logger.Println("Received text document completion uri: " + document.URI)
	urls, err := parsers.ExtractUrls(s.Documents[document.URI])
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
	ignores, err := parsers.ExtractIgnores(s.Documents[document.URI])
	var selectors []master.Selector
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
			Detail:        "from: " + selector.Url,
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

// getSelectors gets all the selectors from the given URL and appends them to the selectors slice
func (s State) getSelectors(
	url string,
	ignores []string,
) ([]master.Selector, error) {
	ctx := context.Background()
	sels, err := s.Database.Queries.ListSelectorsByURL(
		ctx,
		master.ListSelectorsByURLParams{Url: url},
	)
	if err != nil || len(sels) == 0 {
		doc, err := minify.GetMinifiedDoc(url, ignores)
		if err != nil {
			s.Logger.Printf("failed to get minified doc: %s\n", err)
		}
		got := parsers.GetAllSelectors(doc)
		var selectors []master.Selector
		for _, sel := range got {
			selectors = append(selectors, master.Selector{
				Selector: sel,
				Url:      url,
			})
			if _, err := s.Database.Queries.InsertSelector(
				ctx,
				master.InsertSelectorParams{
					Selector: sel,
					Url:      url,
				},
			); err != nil {
				s.Logger.Printf("failed to insert selector: %s\n", err)
			}
		}
		return selectors, nil
	}
	res := []master.Selector{}
	for _, val := range sels {
		res = append(res, master.Selector{
			Url:      url,
			Selector: val.Selector,
			Context:  val.Context,
		})
	}
	return res, nil
}
