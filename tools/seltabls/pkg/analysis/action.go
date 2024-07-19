package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
)

// TextDocumentCodeAction returns the code actions for a given text document.
func TextDocumentCodeAction(
	ctx context.Context,
	req lsp.CodeActionRequest,
	s *State,
) (response rpc.MethodActor, err error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			text := s.Documents[req.Params.TextDocument.URI]
			actions := []lsp.CodeAction{}
			for row, line := range strings.Split(text, "\n") {
				idx := strings.Index(line, "VS Code")
				if idx >= 0 {
					replaceChange := map[string][]lsp.TextEdit{}
					replaceChange[req.Params.TextDocument.URI] = []lsp.TextEdit{
						{
							Range: lsp.LineRange(
								row,
								idx,
								idx+len("VS Code"),
							),
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
							Range: lsp.LineRange(
								row,
								idx,
								idx+len("VS Code"),
							),
							NewText: "VS C*de",
						},
					}
					actions = append(actions, lsp.CodeAction{
						Title: "Censor to VS C*de",
						Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
					})
				}
			}
			return &lsp.TextDocumentCodeActionResponse{
				Response: lsp.Response{
					RPC: lsp.RPCVersion,
					ID:  req.ID,
				},
				Result: actions,
			}, nil
		}
	}
}
