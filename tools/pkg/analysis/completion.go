package analysis

import (
	"log"
	"os"

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
		log.Fatal(err)
	}
	return log.New(logFile, "[seltabl-lsp]", log.LstdFlags)
}

// Selector
type Selector struct {
	Name string
	URL  string
}

// TextDocumentCompletion returns the completions for a given text document.
func (s *State) TextDocumentCompletion(
	id int,
	document *lsp.TextDocumentIdentifier,
	location *lsp.Position,
) lsp.CompletionResponse {
	logger := getLogger("./seltabl.log")
	logger.Println("Received text document completion uri: " + document.URI)
	urls, err := parsers.ExtractUrls(s.Documents[document.URI])
	if err != nil {
		logger.Printf("failed to extract urls: %s\n", err)
		return lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: nil,
		}
	}
	ignores, err := parsers.ExtractIgnores(s.Documents[document.URI])

	var selectors []Selector
	for _, url := range urls {
		doc, err := minify.GetMinifiedDoc(url, ignores)
		if err != nil {
			logger.Printf("failed to get minified doc: %s\n", err)
			return lsp.CompletionResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  &id,
				},
				Result: nil,
			}
		}
		sels := parsers.GetAllSelectors(doc)
		for _, sel := range sels {
			selectors = append(selectors, Selector{
				Name: sel,
				URL:  url,
			})
		}
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
			Label:         selector.Name,
			Detail:        "from: " + selector.URL,
			Documentation: "A selector for the " + selector.Name,
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
