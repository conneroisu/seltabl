package analysis

import (
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

// TextDocumentCodeAction returns the code actions for a given text document.
func (s *State) TextDocumentCodeAction(
	req lsp.CodeActionRequest,
) (response lsp.TextDocumentCodeActionResponse, err error) {
	// Should be able to refresh selectors from the database by requesting the url
	text := s.Documents[req.Params.TextDocument.URI]
	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[req.Params.TextDocument.URI] = []lsp.TextEdit{
				{
					Range:   lsp.LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})
			censorChange := map[string][]lsp.TextEdit{}
			censorChange[req.Params.TextDocument.URI] = []lsp.TextEdit{
				{
					Range:   lsp.LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Censor to VS C*de",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}
	response = lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  req.ID,
		},
		Result: actions,
	}
	return response, nil
}
